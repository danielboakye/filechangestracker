Add-Type -AssemblyName System.Windows.Forms

# Create the form
$form = New-Object System.Windows.Forms.Form
$form.Text = "Service Control"
$form.Size = New-Object System.Drawing.Size(300,150)
$form.StartPosition = "CenterScreen"

# Create a label
$label = New-Object System.Windows.Forms.Label
$label.Text = "Manage the 'filechangestracker' service:"
$label.AutoSize = $true
$label.Location = New-Object System.Drawing.Point(10,10)
$form.Controls.Add($label)

# Create the Start button
$startButton = New-Object System.Windows.Forms.Button
$startButton.Text = "Start Service"
$startButton.Location = New-Object System.Drawing.Point(10,50)
$startButton.Size = New-Object System.Drawing.Size(120,30)
$startButton.Add_Click({
    Start-Service -Name "filechangestracker"
    [System.Windows.Forms.MessageBox]::Show("Service started!")
})
$form.Controls.Add($startButton)

# Create the Stop button
$stopButton = New-Object System.Windows.Forms.Button
$stopButton.Text = "Stop Service"
$stopButton.Location = New-Object System.Drawing.Point(150,50)
$stopButton.Size = New-Object System.Drawing.Size(120,30)
$stopButton.Add_Click({
    Stop-Service -Name "filechangestracker"
    [System.Windows.Forms.MessageBox]::Show("Service stopped!")
})
$form.Controls.Add($stopButton)

# Show the form
$form.Topmost = $true
$form.Add_Shown({$form.Activate()})
[void]$form.ShowDialog()
