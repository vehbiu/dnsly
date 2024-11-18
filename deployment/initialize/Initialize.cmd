# List network adapters
netsh interface ipv4 show interfaces

# Set DNS for your active adapter (usually "Wi-Fi" or "Ethernet")
netsh interface ipv4 set dns name="Wi-Fi" static 127.0.0.1

# To verify
ipconfig /all