$wdir = $pwd
$odir = Join-Path -Path $pwd -ChildPath "../servercsv"

if (! (Test-Path $odir))
{
    New-Item -Path $odir  -ItemType directory
}

$odir = Convert-Path $odir

# get all excel files in directory
$files = get-childitem $wdir -filter *.xlsx

# get excel application COM object
$excel = new-object -comobject excel.application
$excel.visible = $false
$excel.displayalerts = $false

foreach ($file in $files)
{
	$excelfile = $file.fullname
	$filename = $file.basename
	$outfileName = "$odir\" + $filename + ".csv"
	
	Write-Host $excelfile
	$wb = $excel.workbooks.open($excelfile)
	$sheet = $wb.ActiveSheet
	$del = New-Object System.Collections.ArrayList
	for($b = 1 ; $b -lt 10000; $b++)
	{
		$rowvalue = $sheet.cells.item(2, $b).Text
		if ($rowValue -eq "") {
			break
		}
		
		$tmpvalue = $rowvalue -split ","
		
		if ($tmpvalue[0] -eq "client") {
			$del.Add($b)
			continue
		}
		
		if ($tmpvalue[0] -eq "none") {
			$del.Add($b)
			continue
		}
	}
	
	#delete
	$delta = 0
	foreach($i in $del) {
		[void]$sheet.cells.Item(2, $i-$delta).EntireColumn.Delete()
		$delta++
	}
	
	foreach($ws in $wb)
	{
		$ws.saveas($outfileName, 6)
		break
	}

	$wb.close()
	
	#to-utf8
	(gc $outfileName) | Out-File -Encoding 'UTF8' $outfileName
}

# kill excel application
$excel.Quit()
[System.Runtime.Interopservices.Marshal]::ReleaseComObject($excel)
Remove-Variable excel
