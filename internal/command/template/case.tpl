<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>HTTP Request and Response Display</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        h2 { color: #333; }
        .container { margin-bottom: 20px; }
        .section-title { font-weight: bold; margin-top: 20px; }
        pre { background-color: #f4f4f4; padding: 10px; border: 1px solid #ddd; overflow: auto; white-space: pre-wrap; }
    </style>
</head>
<body>
<h2>HTTP Request and Response Display</h2>

<div class="container">
    <div class="section-title">~~~ REQUEST ~~~</div>
    <pre>
    {{ .reqDetail }}
    </pre>
</div>

<div class="container">
    <div class="section-title">~~~ RESPONSE ~~~</div>
    {{if .Error}}

    {{else}}
        <pre>
        {{ .respDetail }}
        </pre>
    {{end}}
</div>
</body>
</html>