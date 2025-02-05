package hash

import (
	"context"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestCreateHashes(t *testing.T) {
	stringsList := []string{"hello", "world"}
	stringsListRes := []string{
		"3338be694f50c5f338814986cdf0686453a888b84f424d792af4b9202398f392",
		"420baf620e3fcd9b3715b42b92506e9304d56e02d3a103499a3a292560cb66b2",
	}

	emptyStringsList := []string{"", " "}
	emptyStringsListRes := []string{
		"a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a",
		"60e893e6d54d8526e55a81f98bfac5da236bb203e84ed5967a8f527d5bf3d4a4",
	}

	intStringsList := []string{"2", "-3"}
	intStringsListRes := []string{
		"b1b1bd1ed240b1496c81ccf19ceccf2af6fd24fac10ae42023628abbe2687310",
		"7ba1e963c14fbce4ad2f5bfe3167be816163c4b29a43925f3f276cbae8296598",
	}

	cyrillicChineseStringsList := []string{"Привет", "世界"}
	cyrillicChineseStringsListRes := []string{
		"f4b999296ad484b981dcfc63ccc275db65ce299d3b55fbf2dace6ad9eb5998d1",
		"aeef5e00d8e45d67c3bca26fb2ffd79ffa68d6ae69fe093fca675adc94348c92",
	}

	emojiList := []string{"😊", "🔥🔥🔥"}
	emojiListRes := []string{
		"bcc7f78f1b69f99a717f5d9bc259fc7590c263004d75b6fdfd8c36edf58cc699",
		"c8bd75d98379ba9f74f8c9069b16c4a9a3b5c08fd3c723286c43f4bb62cd309d",
	}

	longStringList := []string{strings.Repeat("a", 1000), strings.Repeat("b", 2000)}
	longStringListRes := []string{
		"8f3934e6f7a15698fe0f396b95d8c4440929a8fa6eae140171c068b4549fbf81",
		"27399e7506a61a309020fd4d3eaf28b5b678d142ed76e07da117c9114fcc3570",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := map[string]struct {
		args     []string
		expected []string
	}{
		"usual words":          {stringsList, stringsListRes},
		"empty words":          {emptyStringsList, emptyStringsListRes},
		"integer strings":      {intStringsList, intStringsListRes},
		"cyrillic and chinese": {cyrillicChineseStringsList, cyrillicChineseStringsListRes},
		"emoji":                {emojiList, emojiListRes},
		"long strings":         {longStringList, longStringListRes},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			hashes, err := sendToGRPC(ctx, tc.args)

			require.NoError(t, err)
			require.Len(t, hashes, len(stringsList))

			for i, hash := range hashes {
				require.Equal(t, tc.expected[i], hash)
			}
		})
	}
}
