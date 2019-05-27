package libgobuster

import (
	"bytes"
	"fmt"
	"log"

	"github.com/google/uuid"
)

// GobusterDir is the main type to implement the interface
type GobusterDir struct{}

// Setup is the setup implementation of gobusterdir
func (d GobusterDir) Setup(g *Gobuster) error {
	_, _, err := g.GetRequest(g.Opts.URL)
	if err != nil {
		return fmt.Errorf("unable to connect to %s: %v", g.Opts.URL, err)
	}

	guid := uuid.New()
	url := fmt.Sprintf("%s%s", g.Opts.URL, guid)
	wildcardResp, _, err := g.GetRequest(url)
	if err != nil {
		return err
	}

	if g.Opts.StatusCodesParsed.Contains(*wildcardResp) {
		g.IsWildcard = true
		log.Printf("[-] Wildcard response found: %s => %d", url, *wildcardResp)
		if !g.Opts.WildcardForced {
			return fmt.Errorf("To force processing of Wildcard responses, specify the '-fw' switch.")
		}
	}

	return nil
}

// Process is the process implementation of gobusterdir
func (d GobusterDir) Process(g *Gobuster, word string) ([]Result, error) {
	suffix := ""
	if g.Opts.UseSlash {
		suffix = "/"
	}

	// Try the DIR first
	url := fmt.Sprintf("%s%s%s", g.Opts.URL, word, suffix)
	dirResp, dirSize, err := g.GetRequest(url)
	if err != nil {
		return nil, err
	}
	var ret []Result
	if dirResp != nil {
		ret = append(ret, Result{
			Entity: fmt.Sprintf("%s%s", word, suffix),
			Status: *dirResp,
			Size:   dirSize,
		})
	}

	// Follow up with files using each ext.
	for ext := range g.Opts.ExtensionsParsed.Set {
		file := fmt.Sprintf("%s.%s", word, ext)
		url = fmt.Sprintf("%s%s", g.Opts.URL, file)
		fileResp, fileSize, err := g.GetRequest(url)
		if err != nil {
			return nil, err
		}

		if fileResp != nil {
			ret = append(ret, Result{
				Entity: file,
				Status: *fileResp,
				Size:   fileSize,
			})
		}
	}
	return ret, nil
}

// ResultToString is the to string implementation of gobusterdir
func (d GobusterDir) ResultToString(g *Gobuster, r *Result) (*string, error) {
	buf := &bytes.Buffer{}

	// Prefix if we're in verbose mode
	if g.Opts.Verbose {
		if g.Opts.StatusCodesParsed.Contains(r.Status) {
			if _, err := fmt.Fprintf(buf, "Found: "); err != nil {
				return nil, err
			}
		} else {
			if _, err := fmt.Fprintf(buf, "Missed: "); err != nil {
				return nil, err
			}
		}
	}

	if g.Opts.StatusCodesParsed.Contains(r.Status) || g.Opts.Verbose {
		if g.Opts.Expanded {
			if _, err := fmt.Fprintf(buf, g.Opts.URL); err != nil {
				return nil, err
			}
		} else {
			if _, err := fmt.Fprintf(buf, "/"); err != nil {
				return nil, err
			}
		}
		if _, err := fmt.Fprintf(buf, r.Entity); err != nil {
			return nil, err
		}

		if !g.Opts.NoStatus {
			if _, err := fmt.Fprintf(buf, " (Status: %d)", r.Status); err != nil {
				return nil, err
			}
		}

		if r.Size != nil {
			if _, err := fmt.Fprintf(buf, " [Size: %d]", *r.Size); err != nil {
				return nil, err
			}
		}
		if _, err := fmt.Fprintf(buf, "\n"); err != nil {
			return nil, err
		}
	}
	s := buf.String()
	return &s, nil
}
