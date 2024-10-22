<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Styled Table Example</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f4;
            padding: 20px;
        }
        table {
            width: 80%;
            margin: 0 auto;
            border-collapse: collapse;
            background-color: white;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        }
        th, td {
            padding: 12px 15px;
            text-align: left;
        }
        th {
            background-color: #007BFF;
            color: white;
        }
        tr:nth-child(even) {
            background-color: #f2f2f2;
        }
        a {
            color: #007BFF;
            text-decoration: none;
        }
        a:hover {
            text-decoration: underline;
        }
        caption {
            font-size: 1.5em;
            margin: 10px;
            font-weight: bold;
        }
    </style>
</head>
<body>

<table>
    <caption>FinishCount: {{.info.Total}}, SuccessCount: {{.info.SuccessCount}}, FailedCount: {{.info.FailedCount}}</caption>
    <thead>
    <tr>
        <th>id</th>
        <th>description</th>
        <th>state</th>
        <th>reason</th>
        <th>content</th>
    </tr>
    </thead>
    <tbody>
    {{range .tcResultList}}
        <tr>
            <td>{{ .ID }}</td>
            <td>{{ .Desc }}</td>
            <td>{{ .State.String() }}</td>
            <td>{{ .Reason.String() }}</td>
            <td><a href="{{.dirName}}/{{ .ID }}.html">View Details</a></td>
        </tr>
    {{end}}
    </tbody>
</table>

</body>
</html>