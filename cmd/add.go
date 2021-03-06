package cmd

import (
	"fmt"
	"github.com/nicolasmanic/tefter/model"
	"github.com/nicolasmanic/tefter/repository"
	"github.com/spf13/cobra"
	"log"
)

var addNoteCmd = &cobra.Command{
	Use:   "add",
	Short: "Create a new note",
	Long: "A note consist of 4 parts:" +
		" 1) Title, is set through -t flag (optional) \n" +
		" 2) Tags, is set through --tags flag (optional) \n" +
		" 3) Notebook title, if notebook does not exist it will be created,\n" +
		"    is set through -n flag (optional), if not set note will be inserted to the default notebook \n" +
		" 4) Memo, is inserted via VI editor\n",
	Example: "add -t title_1 --tags tag1,tag2 -n notebook_1",
	Run:     addWrapper,
}

func init() {
	rootCmd.AddCommand(addNoteCmd)
	addNoteCmd.Flags().StringP("title", "t", "", "Notes title.")
	addNoteCmd.Flags().StringSlice("tags", []string{}, "Comma-separated tags of note.")
	addNoteCmd.Flags().StringP("notebook", "n", "", "Notebook that this note belongs to")
}

func addWrapper(cmd *cobra.Command, args []string) {
	title, _ := cmd.Flags().GetString("title")
	tags, _ := cmd.Flags().GetStringSlice("tags")
	notebookTitle, _ := cmd.Flags().GetString("notebook")
	editor := &viEditor{}
	err := add(title, tags, notebookTitle, editor)
	if err != nil {
		log.Fatalln(err)
	}
}

func add(title string, tags []string, notebookTitle string, editor Editor) error {
	memo := editor.edit("")

	jNote := &jsonNote{
		Title:         title,
		Memo:          memo,
		Tags:          tags,
		NotebookTitle: notebookTitle,
	}
	return addJSONNote(jNote)
}

func addJSONNote(jNote *jsonNote) error {
	//All newNotes will be inserted to default notebook
	//In next steps the notebook may change see addNotebookToNote for more.
	note := model.NewNote(jNote.Title, jNote.Memo, repository.DEFAULT_NOTEBOOK_ID, jNote.Tags)
	err := addNotebookToNote(note, jNote.NotebookTitle)
	if err != nil {
		return fmt.Errorf("Error while finding corresponding notebook for note, error msg: %v", err)
	}

	_, err = NoteDB.SaveNote(note)
	if err != nil {
		return fmt.Errorf("Error while saving note, error msg: %v", err)
	}
	return nil
}

//addNotebookToNote finds the corresponting notebook for given notebook title
//If notebookTitle exists it will be inserted there.
//If notebookTitle is empty it will be inserted to the default notebook.
//If notebookTitle does not exists notebook will be created and note will be there.
func addNotebookToNote(note *model.Note, notebookTitle string) error {
	if notebookTitle == "" {
		note.NotebookID = repository.DEFAULT_NOTEBOOK_ID
		return nil
	}

	notebook, err := NotebookDB.GetNotebookByTitle(notebookTitle)
	if err != nil {
		return err
	}

	if notebook == nil {
		newNotebook := model.NewNotebook(notebookTitle)
		id, err := NotebookDB.SaveNotebook(newNotebook)
		if err != nil {
			return err
		}
		note.UpdateNotebook(id)
	} else {
		note.UpdateNotebook(notebook.ID)
	}
	return nil
}
