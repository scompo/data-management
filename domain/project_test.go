package domain

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	require := require.New(t)
	p := Project{
		Name:         "test project",
		Description:  "test project description",
		CreationDate: time.Now(),
		LastEdit:     time.Now(),
	}
	ps := []Project{p}
	repo := ProjectsRepository{
		projects: ps,
	}
	repo.Add(ps)
	actual := repo.projects[0]
	require.Equal(p.Name, actual.Name)
	require.Equal(p.Description, actual.Description)
	require.Equal(p.CreationDate, actual.CreationDate)
	require.Equal(p.LastEdit, actual.LastEdit)
}

func TestAll(t *testing.T) {
	require := require.New(t)
	p := Project{
		Name:         "test project",
		Description:  "test project description",
		CreationDate: time.Now(),
		LastEdit:     time.Now(),
	}
	ps := []Project{p}
	repo := ProjectsRepository{
		projects: ps,
	}
	actual := repo.All()
	require.Equal(1, len(actual), "not all data returned")
	require.Equal(p.Name, actual[0].Name)
	require.Equal(p.Description, actual[0].Description)
	require.Equal(p.CreationDate, actual[0].CreationDate)
	require.Equal(p.LastEdit, actual[0].LastEdit)
}
