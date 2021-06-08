package templates

// CliRunner is the outut for a
// job run from the CLI.
const CliRunner = `
{{.banner}}
Credentials:
  - Username:       {{.params.Credentials.Username}}
  - Password:       **********
  - SSH Key File:   {{.params.Credentials.SSHKeyFile}}
  - Super Password: {{.params.Credentials.SuperPassword}}

Devices:
{{- range .params.Devices.Devices }}
  - Name:      {{.Name}}
    IP:        {{.IP}}
    Vendor:    {{.Vendor}}
    Platform:  {{.Platform}}
    Connector: {{.Connector}}
{{- end }}

Commands:
{{- range .params.Commands.Commands}}
  - {{.}}
{{- end }}
{{/* SPACE */}}
`

// CliResult is used to display
// the result of a job run
const CliResult = `{{/* SPACE */}}
{{.Device}}:
  OK: {{.OK}}
  Error: {{.Error}}
  Timestamp: {{.Timestamp}}
`
