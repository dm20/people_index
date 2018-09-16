/*
* AUTHOR: Daniel McGrath
* DATE: 9/13/18
* DESCRIPTION:
* This service tracks "people" information by creating a JSON file called 
* "people.json" in the location that the service is stored, if one does not already exist, 
* and can add new people as well as retrieve people info, and delete people info.
*/

package people_index

import (
  "fmt"
  "os"
  "strings"
  "io/ioutil"
  "github.com/Jeffail/gabs"
)

var existing_json *gabs.Container
var num_deletions int

// func main() {
//   initialize()
//   // call runTests() here 
//   addNewPerson("a key","a value")
//   exit()
// }

// create a new JSON object for each person that is created
func AddNewPerson(key string, value string) {
  existing_json.Set(value,"people",key)
}

// add another value to the given key (creates an array)
func AddChildToKey(key string, value string) {
  newChild := gabs.New()
  newChild.Set(value,"people",key)
  existing_json.Merge(newChild);
}

// list the associated value (single person or array) for a given key
func ListValueForKey(key string) string {
  return existing_json.Path("people." + key).String()
}

// delete a person associated with a given key
// see clearHangingKeys() as well
func DeletePerson(key string) {
  num_deletions++
  existing_json.Set("nil","people",key)
}

// delete any lines in the JSON file marked as deleted
// called whenever a session ends
func ClearHangingKeys(input string) string {
  if (num_deletions == 0) { return input }
  lines := strings.Split(string(input), "\n")
  newLines := make([]string, len(lines) - num_deletions)
  j := 0
  for _, line := range lines {
    if (!strings.Contains(line, "nil")) { 
      newLines[j] = line
      j++
    }
  }
  num_deletions = 0 // reset if more than one call per session
  return strings.Join(newLines, "\n")
}

// Update and close the file. 
func Exit() {
  str := existing_json.StringIndent("", "  ")
  newFile := clearHangingKeys(str) // Any entries deleted in this session are removed
  f, _ := os.Create("./people.json")
  f.Sync()
  f.WriteString(newFile);
  defer f.Close()
}

// If no file exists, create it, otherwise open it in writable format
func Initialize() {
  f, err := os.OpenFile("./people.json",os.O_APPEND | os.O_WRONLY, 0600)
  if (err != nil) {
    fmt.Println("\n\nCreating 'people.json'...\n\n")
    f, _ = os.Create("./people.json")
    jsonObj, _ := gabs.ParseJSON([]byte(`{"people":{}}`))
    newFile := jsonObj.StringIndent("", "  ")
    f.WriteString(newFile)
    fmt.Println("\n\nDone.\n\n")
  }

  existing_data, _ := ioutil.ReadFile("./people.json")
  existing_json, _ = gabs.ParseJSON(existing_data)
}

func RunTests() {
  addNewPerson("1","Oliver")
  addChildToKey("1","A similar person to Oliver, Bill")
  addNewPerson("5","Eric")
  addChildToKey("5","A similar person to Eric, Jane")

  addNewPerson("2","John")
  addNewPerson("3","Amy")
  addNewPerson("4","Reece")
  addNewPerson("6","Allice")
  addNewPerson("7","Bob")
  addNewPerson("8","Joanne")
  addNewPerson("9","Dani")

  deletePerson("3")
  deletePerson("9")
  deletePerson("8")
  deletePerson("7")

  name1 := listValueForKey("1")
  name2 := listValueForKey("2")
  name3 := listValueForKey("3") // will print "nil" and hanging key is deleted when session ends
  name4 := listValueForKey("4")
  name5 := listValueForKey("5")
  name6 := listValueForKey("6")
  fmt.Println(name1)
  fmt.Println(name2)
  fmt.Println(name3)
  fmt.Println(name4)
  fmt.Println(name5)
  fmt.Println(name6)
}