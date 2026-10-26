param(
    [switch]$Install
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

function Test-Tool {
    param([string]$Name)
    return [bool](Get-Command $Name -ErrorAction SilentlyContinue)
}

function Test-Administrator {
    $identity = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($identity)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

function Get-UbuntuState {
    if (-not (Test-Tool 'wsl.exe')) { return '缺失' }
    $previousPreference = $ErrorActionPreference
    $ErrorActionPreference = 'Continue'
    try {
        $listing = (& wsl.exe --list --verbose 2>$null | Out-String) -replace "`0", ''
        $exitCode = $LASTEXITCODE
    } finally {
        $ErrorActionPreference = $previousPreference
    }
    if ($exitCode -ne 0 -or $listing -notmatch 'Ubuntu') { return '缺失' }
    if ($listing -match '(?m)\bUbuntu(?:-[0-9.]+)?\s+\S+\s+2\s*$') { return 'WSL2' }
    return '已安装，但不是 WSL2 或版本无法识别'
}

function Get-DockerState {
    if (-not (Test-Tool 'docker.exe')) { return '缺失' }
    $previousPreference = $ErrorActionPreference
    $ErrorActionPreference = 'Continue'
    try {
        $engine = & docker.exe info --format '{{.OSType}}' 2>$null
        $exitCode = $LASTEXITCODE
    } finally {
        $ErrorActionPreference = $previousPreference
    }
    if ($exitCode -ne 0) { return '已安装，Linux 引擎未运行' }
    if (($engine | Out-String).Trim() -ne 'linux') { return '已安装，当前不是 Linux 容器模式' }
    return '可用'
}

function Show-Status {
    $ubuntu = Get-UbuntuState
    $docker = Get-DockerState
    $memoryGb = [math]::Round((Get-CimInstance Win32_ComputerSystem).TotalPhysicalMemory / 1GB, 1)
    $systemDrive = Get-PSDrive -Name ($env:SystemDrive.TrimEnd(':'))
    $freeGb = [math]::Round($systemDrive.Free / 1GB, 1)
    $processor = Get-CimInstance Win32_Processor | Select-Object -First 1
    $virtualization = if ($processor.VirtualizationFirmwareEnabled) { '已开启' } else { '未开启或无法确认' }
    $items = @(
        [pscustomobject]@{ 项目 = 'Windows 11'; 状态 = '可用' },
        [pscustomobject]@{ 项目 = '内存'; 状态 = ('{0} GB，建议至少 16 GB' -f $memoryGb) },
        [pscustomobject]@{ 项目 = '系统盘剩余空间'; 状态 = ('{0} GB，建议至少留 20 GB' -f $freeGb) },
        [pscustomobject]@{ 项目 = '硬件虚拟化'; 状态 = $virtualization },
        [pscustomobject]@{ 项目 = 'WinGet'; 状态 = $(if (Test-Tool 'winget.exe') { '可用' } else { '缺失' }) },
        [pscustomobject]@{ 项目 = 'WSL2 Ubuntu'; 状态 = $ubuntu },
        [pscustomobject]@{ 项目 = 'Git'; 状态 = $(if (Test-Tool 'git.exe') { '可用' } else { '缺失' }) },
        [pscustomobject]@{ 项目 = 'Go'; 状态 = $(if (Test-Tool 'go.exe') { '可用' } else { '缺失' }) },
        [pscustomobject]@{ 项目 = 'Docker Linux 引擎'; 状态 = $docker },
        [pscustomobject]@{ 项目 = 'kubectl'; 状态 = $(if (Test-Tool 'kubectl.exe') { '可用' } else { '缺失' }) },
        [pscustomobject]@{ 项目 = 'kind'; 状态 = $(if (Test-Tool 'kind.exe') { '可用' } else { '缺失' }) }
    )
    $items | Format-Table -AutoSize | Out-Host
    return [pscustomobject]@{ Ubuntu = $ubuntu; Docker = $docker; MemoryGb = $memoryGb; FreeGb = $freeGb; Virtualization = $virtualization }
}

$os = Get-CimInstance Win32_OperatingSystem
if ([int]$os.BuildNumber -lt 22000 -or $os.Caption -notmatch 'Windows 11') {
    Write-Error '此脚本仅支持 Windows 11。请先确认家里电脑的系统版本。'
    exit 2
}

Write-Host 'AI 应用后端学习环境检查'
$state = Show-Status
if (-not $Install) {
    Write-Host '只进行了检查，没有安装或修改任何内容。'
    exit 0
}

if ($state.MemoryGb -lt 8) {
    Write-Error '内存不足 8 GB，不满足 Docker Desktop 的最低要求；请先确认硬件条件。'
    exit 7
}
if ($state.FreeGb -lt 20) {
    Write-Warning '系统盘剩余空间不足 20 GB；继续安装可能导致 Docker 镜像或本地集群空间不足。'
}
if ($state.Ubuntu -eq '缺失' -and $state.Virtualization -ne '已开启') {
    Write-Error '硬件虚拟化未开启或无法确认。请先在 BIOS/UEFI 中确认虚拟化已开启，再安装 WSL2。'
    exit 8
}

if (-not (Test-Tool 'winget.exe')) {
    Write-Error '缺少 WinGet。请从 Microsoft Store 安装“应用安装程序”，然后重新运行环境准备.cmd。'
    exit 3
}

if ($state.Ubuntu -eq '缺失') {
    if (-not (Test-Administrator)) {
        Write-Host '安装 WSL2 Ubuntu 需要管理员权限，正在请求 Windows 授权。'
        $arguments = '-NoProfile -ExecutionPolicy Bypass -File "{0}" -Install' -f $PSCommandPath
        try {
            $child = Start-Process -FilePath 'powershell.exe' -ArgumentList $arguments -Verb RunAs -Wait -PassThru
            exit $child.ExitCode
        } catch {
            Write-Error '未获得管理员权限，WSL2 Ubuntu 未安装。'
            exit 4
        }
    }
    Write-Host '正在通过 Windows 官方命令安装 WSL2 Ubuntu。'
    & wsl.exe --install -d Ubuntu
    if ($LASTEXITCODE -ne 0) { Write-Error 'WSL2 Ubuntu 安装失败，请查看上面的错误。'; exit 5 }
    Write-Host '请按系统提示重启并完成 Ubuntu 首次用户设置，之后再运行一次环境准备.cmd。'
    exit 0
}

if ($state.Ubuntu -ne 'WSL2') {
    Write-Error 'Ubuntu 已存在，但未确认运行在 WSL2。请先检查 wsl --list --verbose；此脚本不会自动转换已有发行版。'
    exit 6
}

$packages = @(
    [pscustomobject]@{ Tool = 'git.exe'; Id = 'Git.Git' },
    [pscustomobject]@{ Tool = 'go.exe'; Id = 'GoLang.Go' },
    [pscustomobject]@{ Tool = 'docker.exe'; Id = 'Docker.DockerDesktop' },
    [pscustomobject]@{ Tool = 'kubectl.exe'; Id = 'Kubernetes.kubectl' },
    [pscustomobject]@{ Tool = 'kind.exe'; Id = 'Kubernetes.kind' }
)

$failed = @()
foreach ($package in $packages) {
    if (Test-Tool $package.Tool) { continue }
    Write-Host ('正在安装 {0} ...' -f $package.Id)
    & winget.exe install --id $package.Id --exact --source winget
    if ($LASTEXITCODE -ne 0) { $failed += $package.Id }
}

Write-Host '安装阶段结束。请打开新的终端，再运行脚本检查状态。'
if ($failed.Count -gt 0) {
    Write-Warning ('以下安装未成功，请查看上面的 WinGet 输出：' + ($failed -join '、'))
}
Write-Host 'Docker Desktop 首次启动可能需要你完成界面设置；本脚本不代替你接受其协议。'
Write-Host '模型 API 凭证需要自行申请，并以环境变量或本地配置提供；不要写入 Git 仓库。'
