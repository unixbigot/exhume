// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/xml"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var moods = map[int]string{
	90:  "accomplished",
	1:   "aggravated",
	44:  "amused",
	2:   "angry",
	3:   "annoyed",
	4:   "anxious",
	114: "apathetic",
	108: "artistic",
	87:  "awake",
	110: "bitchy",
	92:  "blah",
	113: "blank",
	5:   "bored",
	59:  "bouncy",
	91:  "busy",
	68:  "calm",
	125: "cheerful",
	99:  "chipper",
	84:  "cold",
	63:  "complacent",
	6:   "confused",
	101: "contemplative",
	64:  "content",
	8:   "cranky",
	7:   "crappy",
	106: "crazy",
	107: "creative",
	129: "crushed",
	56:  "curious",
	104: "cynical",
	9:   "depressed",
	45:  "determined",
	130: "devious",
	119: "dirty",
	55:  "disappointed",
	10:  "discontent",
	127: "distressed",
	35:  "ditzy",
	115: "dorky",
	40:  "drained",
	34:  "drunk",
	98:  "ecstatic",
	79:  "embarrassed",
	11:  "energetic",
	12:  "enraged",
	13:  "enthralled",
	80:  "envious",
	78:  "exanimate",
	41:  "excited",
	14:  "exhausted",
	67:  "flirty",
	47:  "frustrated",
	93:  "full",
	103: "geeky",
	120: "giddy",
	72:  "giggly",
	38:  "gloomy",
	126: "good",
	132: "grateful",
	51:  "groggy",
	95:  "grumpy",
	111: "guilty",
	15:  "happy",
	16:  "high",
	43:  "hopeful",
	17:  "horny",
	83:  "hot",
	18:  "hungry",
	52:  "hyper",
	116: "impressed",
	48:  "indescribable",
	65:  "indifferent",
	19:  "infuriated",
	128: "intimidated",
	20:  "irate",
	112: "irritated",
	133: "jealous",
	21:  "jubilant",
	33:  "lazy",
	75:  "lethargic",
	76:  "listless",
	22:  "lonely",
	86:  "loved",
	39:  "melancholy",
	57:  "mellow",
	36:  "mischievous",
	23:  "moody",
	37:  "morose",
	117: "naughty",
	97:  "nauseated",
	102: "nerdy",
	134: "nervous",
	60:  "nostalgic",
	124: "numb",
	61:  "okay",
	70:  "optimistic",
	58:  "peaceful",
	73:  "pensive",
	71:  "pessimistic",
	24:  "pissedoff",
	109: "pleased",
	118: "predatory",
	89:  "productive",
	105: "quixotic",
	77:  "recumbent",
	69:  "refreshed",
	123: "rejected",
	62:  "rejuvenated",
	53:  "relaxed",
	42:  "relieved",
	54:  "restless",
	100: "rushed",
	25:  "sad",
	26:  "satisfied",
	46:  "scared",
	122: "shocked",
	82:  "sick",
	66:  "silly",
	49:  "sleepy",
	27:  "sore",
	28:  "stressed",
	121: "surprised",
	81:  "sympathetic",
	131: "thankful",
	29:  "thirsty",
	30:  "thoughtful",
	31:  "tired",
	32:  "touched",
	74:  "uncomfortable",
	96:  "weird",
	88:  "working",
	85:  "worried",
}

//
// This is a rough mapping from XML to Go of the useful fields in an LJ post
//
type LiveJournalPost struct {
	XMLName          xml.Name `xml:"event"`
	Itemid           int64    `xml:"itemid"`
	Subject          string   `xml:"subject"`
	Eventtime        string   `xml:"eventtime"`
	Event_timestamp  uint64   `xml:"event_timestamp"`
	Url              string   `xml:"url"`
	Current_mood     string   `xml:"current_mood"`
	OptPreformatted  int32    `xml:"opt_preformatted"`
	Current_music    string   `xml:"current_music"`
	Current_location string   `xml:"current_location"`
	Taglist          string   `xml:"props>taglist"`
	Reply_count      int32    `xml:"reply_count"`
	Picture_keyword  string   `xml:"picture_keyword"`
	EventText        string   `xml:"event"`
}

// Cobra command metadata for the post subcommand
var postCmd = &cobra.Command{
	Use:   "post",
	Short: "Create a Hugo post from a LiveJournal Entry",
	Long: `Parse an XML export record from LiveJournal and 
translate it to a Hugo markdown post.

This command accepts a list of filenames representing XML livejournal
exports, and converts them to Markdown files suitable for 
the static website generator [Hugo](http://gohugo.io).

`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, path := range args {
			processPost(path)
		}
	},
}

// Process an LJ post by reading it in and writing out the Hugo equivalent
func processPost(path string) {
	var post LiveJournalPost
	var err error

	err = readJournalPost(path, &post)
	if err != nil {
		log.Fatalf("Error importing %s", path)
	}

	err = writeHugoPost(path, &post)
	if err != nil {
		log.Fatalf("Error exporting %s", path)
	}
}

// Read a LiveJournal XML file (as prdouced by ljdump)
func readJournalPost(path string, post *LiveJournalPost) error {
	log.Printf("Importing LJ post from %v\n", path)

	rawPost, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Unable to read from %v: %v", path, err)
		return err
	}
	//log.Printf("Raw post: %v [%v]\n", path, string(rawPost))

	//log.Printf("Parsing raw post\n")
	err = xml.Unmarshal(rawPost, &post)
	if err != nil {
		log.Printf("Unable to parse XML from %v: %v", path, err)
		return err
	}
	post.EventText = html.UnescapeString(post.EventText)

	//log.Printf("\n\n\nParse result: %+v", post)

	return nil
}

// Write a simple key value pair in TOML format
func writeParam(f *os.File, param string, value interface{}) {
	f.WriteString(fmt.Sprintf("%s = \"%v\"\n", param, value))
}

// Write a key, array pair in TOML format
func writeListParam(f *os.File, param string, values []string) {
	for i, s := range values {
		values[i] = "\"" + s + "\""
	}
	f.WriteString(fmt.Sprintf("%s = [%s]\n", param, strings.Join(values, ",")))
}

// Write a Markdown file equivalent to a given LiveJournal post
func writeHugoPost(pathBase string, post *LiveJournalPost) error {
	path := pathBase + ".md"
	log.Printf("Exporting Hugo post to %v\n", path)

	f, err := os.Create(path)
	if err != nil {
		log.Printf("Unable to open output: %v", err)
		return err
	}
	defer f.Close()

	//
	// Write the Hugo front material
	//
	f.WriteString("+++\n")

	writeParam(f, "title", html.EscapeString(post.Subject))
	tags := strings.Split(post.Taglist, ",")
	if len(post.Current_mood) > 0 {
		tags = append(tags, "mood: "+post.Current_mood)
	}
	if len(tags) > 0 {
		writeListParam(f, "tags", tags)
	}

	if len(post.Picture_keyword) > 0 {
		writeParam(f, "images", fmt.Sprintf("[\"%s.png\"]", post.Picture_keyword))
	}

	when, err := time.Parse("2006-01-02 15:04:05", post.Eventtime)
	if err != nil {
		log.Printf("Unable to parse date [%s]: %v", post.Eventtime, err)
		return err
	}
	writeParam(f, "date", when.Format("2006-01-02 15:04:05"))
	f.WriteString("+++\n\n")

	//
	// Write the post body
	//
	if post.OptPreformatted == 1 {
		f.WriteString("<pre>\n")
	}
	f.WriteString(strings.Replace(post.EventText, "", "", -1))
	if post.OptPreformatted == 1 {
		f.WriteString("</pre>\n")
	}

	return nil
}

func init() {
	RootCmd.AddCommand(postCmd)
	log.SetFlags(log.LstdFlags | log.Llongfile)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// postCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// postCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
