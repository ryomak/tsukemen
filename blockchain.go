package main

import (
)

type BlockchainSession struct{
}

func NewBlockchainSession() *BlockchainSession{
  return new(BlockchainSession)
}

func (b *BlockchainSession) VoteForCandidate(v Vote) error {
	return nil
}
func (b *BlockchainSession) Result() ([]Vote,error) {
	return nil ,nil
}
