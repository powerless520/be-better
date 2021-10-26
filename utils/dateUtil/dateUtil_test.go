package dateUtil

import (
	"testing"
)


func TestIsKid(t *testing.T) {
	t.Log("IdcardIsKid 000000198811110000 is : ", IsKidForIdcard("000000198811110000"))
	t.Log("IdcardIsKid 000000202011110000 is : ", IsKidForIdcard("000000202011110000"))
	t.Log("IdcardIsKid 000000200309010000 is : ", IsKidForIdcard("000000200309010000"))
}