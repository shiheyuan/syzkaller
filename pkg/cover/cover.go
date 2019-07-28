// Copyright 2015 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

// Package cover provides types for working with coverage information (arrays of covered PCs).
package cover

// 复用这部分结构 内核代码raw cover 范围不会超过 uint32大小
// 这个结构将raw cover序列处理为无序集合 这符合当前算法的设计
type Cover map[uint32]struct{}

func (cov *Cover) Merge(raw []uint32) {
	c := *cov
	if c == nil {
		c = make(Cover)
		*cov = c
	}
	for _, pc := range raw {
		c[pc] = struct{}{}
	}
}

func (cov *Cover) MergeCover(c1 Cover) {
	if c1.Empty() {
		return
	}
	s0 := *cov
	if s0 == nil {
		s0 = make(Cover, c1.Len())
		*cov = s0
	}
	for e, p1 := range c1 {
		if _, ok := s0[e]; !ok {
			s0[e] = p1
		}
	}
}

func (cov Cover) Serialize() []uint32 {
	res := make([]uint32, 0, len(cov))
	for pc := range cov {
		res = append(res, pc)
	}
	return res
}

func (cov Cover) Len() int {
	return len(cov)
}

func FromRaw(raw []uint32) Cover {
	if len(raw) == 0 {
		return nil
	}
	s := make(Cover, len(raw))
	for _, e := range raw {
		s[e] = struct{}{}
	}
	return s
}

func FromU64RawUnchecked(raw []uint64) Cover {
	if len(raw) == 0 {
		return nil
	}
	s := make(Cover, len(raw))
	for _, e := range raw {
		s[uint32(e)] = struct{}{}
	}
	return s
}

func (s Cover) Diff(s1 Cover) Cover {
	if s1.Empty() {
		return nil
	}
	var res Cover
	for e, p1 := range s1 {
		if _, ok := s[e]; ok {
			continue
		}
		if res == nil {
			res = make(Cover)
		}
		res[e] = p1
	}
	return res
}

func (s Cover) Intersection(s1 Cover) Cover {
	if s1.Empty() {
		return nil
	}
	res := make(Cover, len(s))
	for e, p := range s {
		if _, ok := s1[e]; ok {
			res[e] = p
		}
	}
	return res
}

func (cov Cover) DiffRaw(raw []uint32) Cover {
	var rec Cover
	for _, e := range raw {
		if _, ok := cov[e]; ok {
			continue
		}
		if rec == nil {
			rec = make(Cover)
		}
		rec[e] = struct{}{}
	}
	return rec
}

func (cov Cover) Empty() bool {
	return len(cov) == 0
}

func RestorePC(pc uint32, base uint32) uint64 {
	return uint64(base)<<32 + uint64(pc)
}

func CovertPC(pc uint64, base uint32) uint32 {
	return uint32(pc - uint64(base))
}
