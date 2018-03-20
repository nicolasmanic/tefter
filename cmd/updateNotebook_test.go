package cmd

import (
	"fmt"
	"github.com/nicolasmanic/tefter/model"
	"github.com/nicolasmanic/tefter/repository"
	"testing"
)

func TestUpdateNotebook(t *testing.T) {
	oldNotebookDB := NotebookDB
	NotebookDB = mockNotebookDBUpdateNotebook{}
	//Restore interface
	defer func() {
		NotebookDB = oldNotebookDB
	}()
	updateNotebook("notebookTitle", "newTitle")
}

func TestUpdateNotebookShouldPanic(t *testing.T) {
	oldNotebookDB := NotebookDB
	NotebookDB = mockNotebookDBUpdateNotebook{}
	defer func() {
		NotebookDB = oldNotebookDB
		if r := recover(); r == nil {
			t.Errorf("Empty arguments should cause the update cmd to panic")
		}
	}()
	updateNotebookWrapper(nil, []string{})
}

type mockNotebookDBUpdateNotebook struct {
	repository.NotebookRepository
}

func (mDB mockNotebookDBUpdateNotebook) UpdateNotebook(notebook *model.Notebook) error {
	if notebook.Title != "newTitle" {
		return fmt.Errorf("Expected Title: %v, got: %v", "newTitle", notebook.Title)
	}
	return nil
}

func (mDB mockNotebookDBUpdateNotebook) GetNotebookByTitle(notebookTitle string) (*model.Notebook, error) {
	notebook := model.NewNotebook(notebookTitle)
	notebook.ID = 1
	note := model.NewNote("testTitle", "testMemo", notebook.ID, []string{})
	note.ID = 3
	notebook.AddNote(note)
	return notebook, nil
}