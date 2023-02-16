package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"openreplay/backend/pkg/db/types"
)

type CH interface {
	InsertAutocomplete(session *types.Session, msgType, msgValue string) error
}

// Conn contains batches, bulks and cache for all sessions
type Conn struct {
	c       Pool
	batches *BatchSet
	bulks   *BulkSet
	chConn  CH
}

func (conn *Conn) SetClickHouse(ch CH) {
	conn.chConn = ch
}

func NewConn(url string, queueLimit, sizeLimit int) *Conn {
	c, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		log.Fatalf("pgxpool.Connect err: %s", err)
	}
	conn := &Conn{}
	conn.c, err = NewPool(c)
	if err != nil {
		log.Fatalf("can't create new pool wrapper: %s", err)
	}
	conn.bulks = NewBulkSet(conn.c)
	conn.batches = NewBatchSet(conn.c, queueLimit, sizeLimit)
	return conn
}

func (conn *Conn) Close() error {
	conn.bulks.Stop()
	conn.batches.Stop()
	conn.c.Close()
	return nil
}

func (conn *Conn) insertAutocompleteValue(sessionID uint64, projectID uint32, tp string, value string) {
	if len(value) == 0 {
		return
	}
	if err := conn.bulks.Get("autocompletes").Append(value, tp, projectID); err != nil {
		log.Printf("autocomplete bulk err: %s", err)
	}
	if conn.chConn == nil {
		return
	}
	// Send autocomplete data to clickhouse
	if err := conn.chConn.InsertAutocomplete(&types.Session{SessionID: sessionID, ProjectID: projectID}, tp, value); err != nil {
		log.Printf("click house autocomplete err: %s", err)
	}
}

func (conn *Conn) batchQueue(sessionID uint64, sql string, args ...interface{}) {
	conn.batches.batchQueue(sessionID, sql, args...)
}

func (conn *Conn) updateSessionEvents(sessionID uint64, events, pages int) {
	conn.batches.updateSessionEvents(sessionID, events, pages)
}

func (conn *Conn) updateSessionIssues(sessionID uint64, errors, issueScore int) {
	conn.batches.updateSessionIssues(sessionID, errors, issueScore)
}

func (conn *Conn) Commit() {
	conn.bulks.Send()
	conn.batches.Commit()
}

func (conn *Conn) updateBatchSize(sessionID uint64, reqSize int) {
	conn.batches.updateBatchSize(sessionID, reqSize)
}
