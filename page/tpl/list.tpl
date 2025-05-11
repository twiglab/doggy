<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>hahah</title>
	</head>
	<body>
		<table border="1">
			<thead>
				<tr>
					<th>SN</th>
					<th>IpAddr</th>
					<th>UUID</th>
					<th>DeviceID</th>
					<th>LastTime</th>
					<th>Operation</th>
				</tr>
			</thead>
			<tbody>
				{{range .Devices}}
				<tr>
					<td>{{.SN }}</td>
					<td>{{.IpAddr }}</td>
					<td>{{.UUID1 }}</td>
					<td>{{.Code1 }}</td>
					<td>{{.Last }}</td>
					<td><a href="">open</a></td>
				</tr>
				{{end}}
			</tbody>
		</table>
	</body>
</html>
