/*
 * owobot - Your server's guardian and entertainer
 * Copyright (C) 2023 owobot Contributors
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package limiter

import (
	"sync"
	"time"
)

// Limiter is a token bucket rate limiter
type Limiter struct {
	TotalAmt int
	WarnAmt  int
	Duration time.Duration

	mu     sync.Mutex
	tokens map[string]int
}

// New returns a new limiter with the given parameters
func New(warnAmt, totalAmt int, duration time.Duration) *Limiter {
	limiter := &Limiter{
		TotalAmt: totalAmt,
		WarnAmt:  warnAmt,
		Duration: duration,
		tokens:   map[string]int{},
	}
	go limiter.resetTokens()
	return limiter
}

// Decrement removes one token from the bucket with the given key
func (l *Limiter) Decrement(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	tokenAmt, ok := l.tokens[key]
	if ok {
		l.tokens[key] = tokenAmt - 1
	} else {
		l.tokens[key] = l.TotalAmt - 1
	}
}

// IsWarning returns true if the token amount equals the warn amount of l.
func (l *Limiter) IsWarning(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	tokenAmt, ok := l.tokens[key]
	if !ok {
		return false
	}
	return tokenAmt > 0 && tokenAmt <= (l.TotalAmt-l.WarnAmt)
}

// IsWarning returns true if the token amount for the given key is depleted.
func (l *Limiter) IsDepleted(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	tokenAmt, ok := l.tokens[key]
	if !ok {
		return false
	}
	return tokenAmt <= 0
}

// resetTokens resets all the token buckets at a regular interval
func (l *Limiter) resetTokens() {
	for {
		l.mu.Lock()
		l.tokens = map[string]int{}
		l.mu.Unlock()
		time.Sleep(l.Duration)
	}
}
