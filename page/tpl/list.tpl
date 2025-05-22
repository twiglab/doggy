<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Doggy</title>
	</head>
	<body>
		<table border="1">
			<thead>
				<tr>
					<th>SN</th>
					<th>IpAddr</th>
					<th>UUID</th>
					<th>DeviceID</th>
					<th>User</th>
					<th>Pwd</th>
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
					<td>{{.User }}</td>
					<td>{{.Pwd }}</td>
					<td>{{.Last.Format "2006-01-02T15:04:05Z07:00" }}</td>
					<td>
						<a href="https://{{ .IpAddr }}" target="_blank">open</a> |
						Ping | Reboot
					</td>
				</tr>
				{{end}}
			</tbody>
		</table>
	</body>
</html>
