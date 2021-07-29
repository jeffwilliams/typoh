package main

import (
	"bytes"
	"testing"
)

func TestTypographer(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{

		{
			name: "no replacements",
			input: "<html>\n  <body>\n  	This is a ``test''---so they `say'. But don't modify this 20\\' or 6\"...\n  </body>\n</html>\n",
			output: `<html>
  <body>
  	This is a “test”—so they ‘say’. But don’t modify this 20' or 6"…
  </body>
</html>
`,
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			var typo Typographer
			inBuf := bytes.NewBufferString(tc.input)
			var outBuf bytes.Buffer

			typo.ReplaceMarkers(inBuf, &outBuf)
			output := outBuf.String()
			if tc.output != output {
				t.Fatalf("Replacement was wrong for '%s'. Expected %v but got %v", tc.input, tc.output, output)
			}

		})
	}
}
