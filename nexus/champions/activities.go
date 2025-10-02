/*
 * Copyright 2025 Simon Emms <simon@simonemms.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package champions

import (
	"context"
	"math/rand"
	"time"

	"go.temporal.io/sdk/activity"
)

// If you wish to connect any dependencies (eg, database), add in here
type activities struct {
	// Database  *gorm.DB
}

// It's a curious oddity that there have been no repeating given names of any
// F1 champion in 75 years. Yes, Mika and Michael mean the same and it's
// officially Sir John Young Stewart, but he's Jackie (or JYS if you read Jenks).
//
// If Kimi Antonelli becomes champion, he'll be the first. And he's officially
// Andrea Kimi.
var availableNames = []string{
	"Nino",
	"Juan Manuel",
	"Alberto",
	"Mike",
	"Jack",
	"Phil",
	"Graham",
	"Jim",
	"John",
	"Denny",
	"Jackie",
	"Jochen",
	"Emerson",
	"Niki",
	"James",
	"Mario",
	"Jody",
	"Alan",
	"Nelson",
	"Alain",
	"Ayrton",
	"Nigel",
	"Michael",
	"Damon",
	"Jacques",
	"Mika",
	"Fernando",
	"Kimi",
	"Lewis",
	"Jenson",
	"Sebastian",
	"Nico",
	"Max",
	// Who's next?
	"Oscar",
	"Lando",
}

// This is a simple activity that finds a world champion
func (a *activities) FindChampion(ctx context.Context) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Find a champion")

	time.Sleep(time.Second * 5)

	return availableNames[rand.Intn(len(availableNames))], nil
}

// If you need to inject dependencies, pass them in here
// The error response can be useful
func NewActivities( /*dbConnection *gorm.DB*/ ) (*activities, error) {
	return &activities{
		// Database: dbConnection,
	}, nil
}
