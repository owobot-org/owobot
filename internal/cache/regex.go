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

package cache

import (
	"regexp"
	"sync"
)

var (
	regexMtx = sync.RWMutex{}
	regexes  = map[string]*regexp.Regexp{}
)

func Regex(regex string) (*regexp.Regexp, error) {
	regexMtx.RLock()
	if re, ok := regexes[regex]; ok {
		regexMtx.RUnlock()
		return re, nil
	}
	regexMtx.RUnlock()

	re, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}

	regexMtx.Lock()
	regexes[regex] = re
	regexMtx.Unlock()

	return re, nil
}
