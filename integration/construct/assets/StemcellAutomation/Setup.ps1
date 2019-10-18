echo "mock stemcell automation script executed"
Start-Sleep -s 45

if (magic-file-present)
    exit
Stop-Computer