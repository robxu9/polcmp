package main

import "html/template"

const (
	candidateTable = `
  {{range $table := .}}
  <div class="table-responsive">
    <table class="table table-hover">
      <!-- TODO <caption></caption> -->
      <thead>
        <tr>
          <th>{{$table.Name}}</th>
          {{range $table.Candidates}}
          <th>{{.}}</th>
          {{end}}
        </tr>
      </thead>
      <tbody>
        {{range $k, $v := $table.Positions}}
        <tr>
          <th scope="row">{{$k}}</th>
          {{range $v}}
          <td>{{.}}</td>
          {{end}}
        </tr>
				{{end}}
      </tbody>
    </table>
  </div>
  {{end}}
`
)

var (
	tableTmpl = template.Must(template.New("table").Parse(candidateTable))
)

// TABLE NAME | CANDIDATE1 | CANDIDATE2 | CANDIDATE3
// POSITION   | POSITION1  | POSITION2  | POSITION3
// POSITION   | POSITION1  | POSITION2  | POSITION3

// Table represents the table for use on the template.
type Table struct {
	Name       string              // name to put in thead
	Candidates []string            // to map to the columns
	Positions  map[string][]string // position ~> position1,position2,position3
}
