package envcfg

import (
	"reflect"
	"testing"
)

func Test_parse(t *testing.T) {
	tests := []struct {
		name, doc, wantErr string
		want               parsedOpts
	}{
		{name: "empty", doc: "", wantErr: "doc must contain at least name"},
		{
			name: "empty",
			doc: ` 
	`,
			wantErr: "doc must contain at least name",
		},
		{name: "named", doc: "   I_AM_VAR   ", want: parsedOpts{"I_AM_VAR", "", nil}},

		{name: "flag same line", doc: `I_AM_VAR optional`, want: parsedOpts{"I_AM_VAR", "", [][2]string{{"optional"}}}},
		{
			name: "flag next line", doc: `I_AM_VAR
optional`,
			want: parsedOpts{"I_AM_VAR", "", [][2]string{{"optional"}}},
		},
		{
			name: "flag blank lines", doc: `I_AM_VAR

optional`,
			want: parsedOpts{"I_AM_VAR", "", [][2]string{{"optional"}}},
		},

		{
			name: "option same line",
			doc:  `I_AM_VAR default=2`,
			want: parsedOpts{"I_AM_VAR", "", [][2]string{{"default", "2"}}},
		},
		{
			name: "option next line", doc: `I_AM_VAR
default=2`,
			want: parsedOpts{"I_AM_VAR", "", [][2]string{{"default", "2"}}},
		},
		{
			name: "option blank lines", doc: `I_AM_VAR

default=2`,
			want: parsedOpts{"I_AM_VAR", "", [][2]string{{"default", "2"}}},
		},
		{
			name: "option is space",
			doc:  `I_AM_VAR default==`,
			want: parsedOpts{"I_AM_VAR", "", [][2]string{{"default", `=`}}},
		},
		{
			name: "option contains whitespace",
			doc: `I_AM_VAR default=" blah
"`,
			want: parsedOpts{"I_AM_VAR", "", [][2]string{{"default", ` blah
`}}},
		},
		{
			name: "option is quote",
			doc:  `I_AM_VAR default="`,
			want: parsedOpts{"I_AM_VAR", "", [][2]string{{"default", `"`}}},
		},
		{
			name: "option contains quote",
			doc:  `I_AM_VAR default=ab"cd`,
			want: parsedOpts{"I_AM_VAR", "", [][2]string{{"default", `ab"cd`}}},
		},
		{
			name: "option contains multiple quotes",
			doc:  `I_AM_VAR default=ab"cd"ef"`,
			want: parsedOpts{"I_AM_VAR", "", [][2]string{{"default", `ab"cd"ef"`}}},
		},
		{
			name: "option contains escaped quote",
			doc:  `I_AM_VAR default="\"\""`,
			want: parsedOpts{"I_AM_VAR", "", [][2]string{{"default", `""`}}},
		},
		{
			name: "option contains escaped slash",
			doc:  `I_AM_VAR default="\"\\"`,
			want: parsedOpts{"I_AM_VAR", "", [][2]string{{"default", `"\`}}},
		},
		{
			name: "option contains escaped followed by escaped quote",
			doc:  `I_AM_VAR default="\\\""`,
			want: parsedOpts{"I_AM_VAR", "", [][2]string{{"default", `\"`}}},
		},
		{
			name: "option contains elaborate multiline quoted value with delimiter",
			doc: `I_AM_VAR default=" whoa
so \"
frigging'
awesome yeah'\"
"

Optional

---
I_AM_VAR is a super cool weird thing.
`,
			want: parsedOpts{
				"I_AM_VAR",
				"I_AM_VAR is a super cool weird thing.",
				[][2]string{
					{"default", ` whoa
so "
frigging'
awesome yeah'"
`},
					{"Optional"},
				},
			},
		},
		{
			name: "no options just description",
			doc:  `I_AM_VAR I_AM_VAR sets stuff`,
			want: parsedOpts{"I_AM_VAR", "I_AM_VAR sets stuff", nil},
		},
		{
			name: "shorthand form",
			doc:  `I_AM_VAR option default=2 | sets stuff`,
			want: parsedOpts{"I_AM_VAR", "sets stuff", [][2]string{{"option"}, {"default", "2"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse(tt.doc)

			var errS string
			if err != nil {
				errS = err.Error()
			}
			if errS != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
