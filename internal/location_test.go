package internal

import "testing"

func TestGetType(t *testing.T) {

	t.Run("TypeLocal", func(t *testing.T) {
		l := Location{
			Type: "local",
		}
		result, err := l.getType()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		assertEqual(t, result, TypeLocal)
	})

	t.Run("TypeVolume", func(t *testing.T) {
		l := Location{
			Type: "volume",
		}
		result, err := l.getType()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		assertEqual(t, result, TypeVolume)
	})

	t.Run("Empty type", func(t *testing.T) {
		l := Location{
			Type: "",
		}
		result, err := l.getType()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		assertEqual(t, result, TypeLocal)
	})

	t.Run("Invalid type", func(t *testing.T) {
		l := Location{
			Type: "foo",
		}
		_, err := l.getType()
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestBuildTag(t *testing.T) {
	result := buildTag("foo", "bar")
	expected := "ar:foo:bar"
	assertEqual(t, result, expected)
}

func TestGetLocationTags(t *testing.T) {
	l := Location{
		name: "foo",
	}
	result := l.getLocationTags()
	expected := "ar:location:foo"
	assertEqual(t, result, expected)
}

func TestHasBackend(t *testing.T) {
	t.Run("backend present", func(t *testing.T) {
		l := Location{
			name: "foo",
			To:   []string{"foo", "bar"},
		}
		result := l.hasBackend("foo")
		assertEqual(t, result, true)
	})

	t.Run("backend absent", func(t *testing.T) {
		l := Location{
			name: "foo",
			To:   []string{"bar", "baz"},
		}
		result := l.hasBackend("foo")
		assertEqual(t, result, false)
	})
}

func TestBuildRestoreCommand(t *testing.T) {
	l := Location{
		name: "foo",
	}
	result := buildRestoreCommand(l, "to", "snapshot", []string{"options"})
	expected := []string{"restore", "--target", "to", "--tag", "ar:location:foo", "snapshot", "options"}
	assertSliceEqual(t, result, expected)
}
