/*
 * owobot - The coolest Discord bot ever written
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

package xsync

import "sync"

// KeyedMutex is a mutex implementation using several independent keys
type KeyedMutex struct {
	mutexes sync.Map
}

// Lock locks a mutex with the given key. If a mutex with that key
// doesn't exist, a new one is created.
func (m *KeyedMutex) Lock(key string) {
	value, _ := m.mutexes.LoadOrStore(key, &sync.Mutex{})
	mtx := value.(*sync.Mutex)
	mtx.Lock()
}

// Unlock unlocks a mutex with the given key if it exists.
func (m *KeyedMutex) Unlock(key string) {
	value, ok := m.mutexes.Load(key)
	if ok {
		mtx := value.(*sync.Mutex)
		mtx.Unlock()
	}
}
