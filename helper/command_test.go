package helper

import (
	"fmt"
	"testing"

	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func TestGetCommand(t *testing.T) {
	botUsername := "peminjaman_testing_bot"

	t.Run("return the command", func(t *testing.T) {
		m := fmt.Sprintf("/%s", types.CommandBorrow)
		r := GetCommand(m)

		assert.Equal(t, types.CommandBorrow, r)
	})

	t.Run("return empty string if no '/' symbol", func(t *testing.T) {
		m := "haztest"
		r := GetCommand(m)

		assert.Empty(t, r)
	})

	t.Run("return empty string if '/' symbol is not in the first char", func(t *testing.T) {
		m := "haztest/"
		r := GetCommand(m)

		assert.Empty(t, r)
	})

	t.Run("return the command, separated with space char", func(t *testing.T) {
		m := fmt.Sprintf("/%s haztest", types.CommandBorrow)
		r := GetCommand(m)

		assert.Equal(t, types.CommandBorrow, r)
	})

	t.Run("mention", func(t *testing.T) {
		m := fmt.Sprintf("/%s@%s", types.CommandCheck, botUsername)
		r := GetCommand(m)

		assert.Equal(t, types.CommandCheck, r)
	})

	t.Run("mention with another data", func(t *testing.T) {
		m := fmt.Sprintf("/%s@%s %s %d", types.CommandRespond, botUsername, types.RespondTypeBorrow, 123)
		r := GetCommand(m)

		assert.Equal(t, types.CommandRespond, r)
	})
}

func TestGetRespondCommands(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		s := fmt.Sprintf("/%s %s %d yes", types.CommandRespond, types.RespondTypeBorrow, 123)
		r, ok := GetRespondCommandOrder(s)

		expected := types.RespondCommandOrder{
			Type: types.RespondTypeBorrow,
			ID:   123,
			Text: "yes",
		}

		assert.True(t, ok)
		assert.Equal(t, expected, r)
	})

	t.Run("length less than 3", func(t *testing.T) {
		s := fmt.Sprintf("/%s %s", types.CommandRespond, types.RespondTypeBorrow)
		r, ok := GetRespondCommandOrder(s)

		assert.False(t, ok)
		assert.Equal(t, types.RespondCommandOrder{}, r)
	})

	t.Run("length equal 3", func(t *testing.T) {
		s := fmt.Sprintf("/%s %s %d", types.CommandRespond, types.RespondTypeBorrow, 123)
		r, ok := GetRespondCommandOrder(s)

		expected := types.RespondCommandOrder{
			Type: types.RespondTypeBorrow,
			ID:   123,
			Text: "",
		}

		assert.True(t, ok)
		assert.Equal(t, expected, r)
	})

	t.Run("exceed length 4", func(t *testing.T) {
		s := fmt.Sprintf("/%s %s %d yes oke nais 123", types.CommandRespond, types.RespondTypeBorrow, 123)
		r, ok := GetRespondCommandOrder(s)

		assert.False(t, ok)
		assert.Equal(t, types.RespondCommandOrder{}, r)
	})

	t.Run("not in category", func(t *testing.T) {
		s := fmt.Sprintf("/%s testnotincategory %d yes", types.CommandRespond, 123)
		r, ok := GetRespondCommandOrder(s)

		assert.False(t, ok)
		assert.Equal(t, types.RespondCommandOrder{}, r)
	})

	t.Run("wrong id", func(t *testing.T) {
		s := fmt.Sprintf("/%s %s testwrongid yes", types.CommandRespond, types.RespondTypeBorrow)
		r, ok := GetRespondCommandOrder(s)

		assert.False(t, ok)
		assert.Equal(t, types.RespondCommandOrder{}, r)
	})
}

func TestIsRespondTypeExists(t *testing.T) {
	t.Run("exist", func(t *testing.T) {
		r := isRespondTypeExists(types.RespondTypeBorrow)
		assert.True(t, r)
	})
	t.Run("nope", func(t *testing.T) {
		r := isRespondTypeExists("testdoesnotexists")
		assert.False(t, r)
	})
}

func TestIsCommandTypeExists(t *testing.T) {
	t.Run("exist", func(t *testing.T) {
		r := isManageTypeExists(types.ManageTypeAdd)
		assert.True(t, r)
	})
	t.Run("nope", func(t *testing.T) {
		r := isManageTypeExists("testdoesnotexists")
		assert.False(t, r)
	})
}

func TestGetManageCommands(t *testing.T) {
	t.Run("full value", func(t *testing.T) {
		s := fmt.Sprintf("/%s %s %d", types.CommandManage, types.ManageTypeEdit, 123)
		r, ok := GetManageCommandOrder(s)

		expected := types.ManageCommandOrder{
			Type: types.ManageTypeEdit,
			ID:   123,
		}

		assert.True(t, ok)
		assert.Equal(t, expected, r)
	})

	t.Run("length less than 2", func(t *testing.T) {
		s := fmt.Sprintf("/%s", types.CommandManage)
		r, ok := GetManageCommandOrder(s)

		assert.False(t, ok)
		assert.Equal(t, types.ManageCommandOrder{}, r)
	})

	t.Run("length equal 2 with correct type", func(t *testing.T) {
		s := fmt.Sprintf("/%s %s", types.CommandManage, types.ManageTypeAdd)
		r, ok := GetManageCommandOrder(s)

		expected := types.ManageCommandOrder{
			Type: types.ManageTypeAdd,
		}
		assert.True(t, ok)
		assert.Equal(t, expected, r)
	})

	t.Run("length equal 2 with incorrect type", func(t *testing.T) {
		s := fmt.Sprintf("/%s %s", types.CommandManage, types.ManageTypeEdit)
		r, ok := GetManageCommandOrder(s)

		assert.False(t, ok)
		assert.Equal(t, types.ManageCommandOrder{}, r)
	})

	t.Run("length exceed 3", func(t *testing.T) {
		s := fmt.Sprintf("/%s %s %d %s", types.CommandManage, types.ManageTypeEdit, 123, "testexceed")
		r, ok := GetManageCommandOrder(s)

		assert.False(t, ok)
		assert.Equal(t, types.ManageCommandOrder{}, r)
	})

	t.Run("wrong type", func(t *testing.T) {
		s := fmt.Sprintf("/%s testwrongtype %d", types.CommandManage, 123)
		r, ok := GetManageCommandOrder(s)

		assert.False(t, ok)
		assert.Equal(t, types.ManageCommandOrder{}, r)
	})

	t.Run("wrong id", func(t *testing.T) {
		s := fmt.Sprintf("/%s %s testwrongid", types.CommandManage, types.ManageTypeEdit)
		r, ok := GetManageCommandOrder(s)

		assert.False(t, ok)
		assert.Equal(t, types.ManageCommandOrder{}, r)
	})
}

func TestGetCheckCommandOrder(t *testing.T) {
	t.Run("full value", func(t *testing.T) {
		s := fmt.Sprintf("/%s %d %s", types.CommandCheck, 123, types.CheckTypePhoto)
		r, ok := GetCheckCommandOrder(s)

		expected := types.CheckCommandOrder{
			ID:   123,
			Text: types.CheckTypePhoto,
		}

		assert.True(t, ok)
		assert.Equal(t, expected, r)
	})

	t.Run("length less than 2", func(t *testing.T) {
		s := fmt.Sprintf("/%s", types.CommandCheck)
		r, ok := GetCheckCommandOrder(s)

		assert.False(t, ok)
		assert.Equal(t, types.CheckCommandOrder{}, r)
	})

	t.Run("length exceed 3", func(t *testing.T) {
		s := fmt.Sprintf("/%s %d %s testexceed", types.CommandCheck, 123, types.CheckTypePhoto)
		r, ok := GetCheckCommandOrder(s)

		assert.False(t, ok)
		assert.Equal(t, types.CheckCommandOrder{}, r)
	})

	t.Run("wrong id", func(t *testing.T) {
		s := fmt.Sprintf("/%s testwrongid %s", types.CommandCheck, types.CheckTypePhoto)
		r, ok := GetCheckCommandOrder(s)

		assert.False(t, ok)
		assert.Equal(t, types.CheckCommandOrder{}, r)
	})
}

func TestGetReportCommands(t *testing.T) {
	t.Run("full value", func(t *testing.T) {
		s := fmt.Sprintf("/%s %s", types.CommandReport, types.ReportTypeBorrow)
		r, ok := GetReportCommandOrder(s)

		expected := types.ReportCommandOrder{
			Type: types.ReportTypeBorrow,
		}

		assert.True(t, ok)
		assert.Equal(t, expected, r)
	})

	t.Run("length not equal 2", func(t *testing.T) {
		s := fmt.Sprintf("/%s", types.CommandReport)
		r, ok := GetReportCommandOrder(s)

		assert.False(t, ok)
		assert.Equal(t, types.ReportCommandOrder{}, r)
	})

	t.Run("wrong type", func(t *testing.T) {
		s := fmt.Sprintf("/%s testwrongtype", types.CommandReport)
		r, ok := GetReportCommandOrder(s)

		assert.False(t, ok)
		assert.Equal(t, types.ReportCommandOrder{}, r)
	})
}
