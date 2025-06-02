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
					<th>Mac</th>
					<th>UUID</th>
					<th>DeviceID</th>
					<th>User</th>
					<th>Pwd</th>
					<th>TTL</th>
					<th>LastTime</th>
					<th>Operation</th>
				</tr>
			</thead>
			<tbody>
				{{range .Devices}}
				<tr>
					<td>{{ .Upload.SN }}</td>
					<td>{{ .Upload.IpAddr }}</td>
					<td>{{ .Data.Mac }}</td>
					<td>{{ .Data.UUID }}</td>
					<td>{{ .Data.DeviceID }}</td>
					<td>{{ .Upload.User }}</td>
					<td>{{ .Upload.Pwd }}</td>
					<td>{{ .TTL }}</td>
					<td>{{ .Upload.Last.Format "2006-01-02T15:04:05Z07:00" }}</td>
					<td>
						<a href="https://{{ .Upload.IpAddr }}" target="_blank">open</a> |
						Ping | Reboot
					</td>
				</tr>
				{{end}}
			</tbody>
		</table>
	</body>
</html>
