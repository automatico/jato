package outputters

// CliRunner is the outut for a
// job run from the CLI.
const CliRunner = `{{/* SPACE */}}
--------------------------
Job Parameters
--------------------------
Username: {{.User.Username}}
Password: *************
Devices:
{{- range .Devices.Devices }}
  - Name:      {{.Name}}
    IP:        {{.IP}}
    Vendor:    {{.Vendor}}
    Platform:  {{.Platform}}
    Connector: {{.Connector}}
{{- end }}
Commands:
{{- range .Commands.Commands}}
  - {{.}}
{{- end }}

--------------------------
Job Result
--------------------------
{{/* SPACE */}}`
