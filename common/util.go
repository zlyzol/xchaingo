package common

import (
)

func SingleFee(fee Fee) Fees {
	return Fees{
		Type:		FlatFee,
		FeeOptions: FeeOptions{
			Average:	fee,
			Fast:		fee,
			Fastest:	fee,
		},
	  }
}