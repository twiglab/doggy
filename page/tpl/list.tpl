<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Camera List</title>
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
					<td>{{.Last.Format "2006-01-02T15:04:05Z07:00" }}</td>
					<td>
						<a href="https://{{ .IpAddr }}" target="_blank">open</a> |
					</td>
				</tr>
				{{end}}
			</tbody>
		</table>
	</body>
</html>
