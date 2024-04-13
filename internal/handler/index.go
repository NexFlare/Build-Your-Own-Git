package handler

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/codecrafters-io/git-starter-go/internal/object"
	"github.com/codecrafters-io/git-starter-go/internal/util"
)


func CommitHandler() string{
	if len(os.Args) < 5 {
		fmt.Fprintf(os.Stderr, "usage: mygit commit-tree <tree_sha> -m <message>\n")
		os.Exit(1)
	}
	treeSha := os.Args[2]
	parent := ""
	message := ""
	time := time.Now()
	commitSeconds := time.Unix()
	timezone := util.GetTimeZone(time)
	authorObj := object.User {
		Email: "harshmagarwal@gmail.com",
		Name: "Harsh Agarwal",
		Seconds: commitSeconds,
		Timezone: timezone,
	}

	commiterObj := object.User {
		Email: "harshmagarwal@gmail.com",
		Name: "Harsh Agarwal",
		Seconds: commitSeconds,
		Timezone: timezone,
	}

	for i:=3; i<len(os.Args); {
		if os.Args[i] == "-m" {
			i++
			message = os.Args[i]
		} else if os.Args[i] == "-p" {
			i++
			parent = os.Args[i]
		}
		i++
	}

	commitContent := strings.Builder{}
	commitContent.WriteString(fmt.Sprintf("tree %s\n", treeSha))
	if len(parent) > 0 {
		commitContent.WriteString(fmt.Sprintf("parent %s\n", parent))
	}
	commitContent.WriteString(fmt.Sprintf("author %s <%s> %d %s\n", authorObj.Name, authorObj.Email, authorObj.Seconds, authorObj.Timezone))
	commitContent.WriteString(fmt.Sprintf("committer %s <%s> %d %s\n", commiterObj.Name, commiterObj.Email, commiterObj.Seconds, commiterObj.Timezone))
	commitContent.WriteString(fmt.Sprintf("\n%s\n", message))
	commitData := fmt.Sprintf("commit %d\x00%s", commitContent.Len(), commitContent.String())
	hashedValue := util.GetHash([]byte(commitData))
	util.CreateObjectFile([]byte(commitData), hex.EncodeToString(hashedValue))
	return hex.EncodeToString(hashedValue)
}