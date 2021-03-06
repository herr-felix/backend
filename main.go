/* Parle - a free customer communication platform
 * Copyright (C) 2017 Focus Centric inc.
 *
 * You may use, distribute and modify this code under the
 * terms of the GNU Affero General Public license, as
 * published by the Free Software Foundation, either version
 * 3 of theLicense, or (at your option) any later version.
 *
 * You should have received a copy of the GNU Affero General
 * Public License along with this code as LICENSE file.  If not,
 * see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/parle-io/backend/data/sqlite"
)

var (
	addr = flag.String("addr", ":8080", "http service address")
	db   = flag.String("db", "db/parle.db", "database path and file name (default to db/dev.db)")
)

func main() {
	flag.Parse()

	dbName := *db
	if len(dbName) == 0 {
		dbName = "db/dev.db"
	}

	// Initiating the database connection pool
	sqlitePersistence := sqlite.Open(dbName)

	hub := newHub(sqlitePersistence)
	go hub.run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("unable to start server: ", err)
	}
}
