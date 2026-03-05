package kafka



type commit struct {
	topic     string
	partition int
	offset    int64
}



func makeCommit(msg Message) commit {
	return commit{
		topic:     msg.Topic,
		partition: msg.Partition,
		offset:    msg.Offset + 1,
	}
}




func makeCommits(msgs ...Message) []commit {
	commits := make([]commit, len(msgs))

	for i, m := range msgs {
		commits[i] = makeCommit(m)
	}

	return commits
}



type commitRequest struct {
	commits []commit
	errch   chan<- error
}
