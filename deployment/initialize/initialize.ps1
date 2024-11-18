# List network adapters first to find the one you want to modify
Get-NetAdapter

# Set primary DNS to localhost for the active network adapter (usually "Ethernet" or "Wi-Fi")
Set-DnsClientServerAddress -InterfaceAlias "Wi-Fi" -ServerAddresses "127.0.0.1"

# To verify the change
Get-DnsClientServerAddress -InterfaceAlias "Wi-Fi"