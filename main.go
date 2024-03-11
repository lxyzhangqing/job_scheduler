package main

import "fmt"

type Job struct {
	Id    int
	Cpu   int
	Mem   int
	Value int
}

var jobs = []Job{
	{Id: 1, Cpu: 1, Mem: 1, Value: 1},
	{Id: 2, Cpu: 1, Mem: 10, Value: 1},
	{Id: 3, Cpu: 1, Mem: 2, Value: 1},
	{Id: 4, Cpu: 2, Mem: 8, Value: 1},
	{Id: 5, Cpu: 2, Mem: 3, Value: 1},
	{Id: 6, Cpu: 2, Mem: 1, Value: 1},
} // 最优解应该输出1、3、4、5、6

var (
	cpuLimit = 8
	memLimit = 16
)

type State struct {
	Value int
	Ids   []int
}

func (s *State) AddId(id int) {
	if s.Ids == nil {
		s.Ids = []int{}
	}
	s.Ids = append(s.Ids, id)
}

func (s State) DeepCopy() State {
	ns := s
	if s.Ids != nil {
		ns.Ids = append([]int{}, s.Ids...)
	}
	return ns
}

var state = [6][9][17]State{}

func main() {
	for i := 0; i < len(jobs); i++ {
		for c := 0; c <= cpuLimit; c++ {
			for m := 0; m <= memLimit; m++ {
				// 状态转移函数
				stateTransition(i, c, m)
			}
		}
	}

	result()
}

func stateTransition(i, c, m int) {
	// 如果job i因资源限额不能放入机器
	if overResourceLimit(jobs[i], c, m) {
		// 如果job i是第1个job，那么只在前i个job中挑选的最优值为0
		if i == 0 {
			state[i][c][m].Value = 0
		} else {
			// 如果job i不是第1个job，那么只在前i个job中挑选的最优值，与只在前i-1个job中挑选的最优值相同
			state[i][c][m] = state[i-1][c][m].DeepCopy()
		}
	} else { //如果job i能放入机器中

		if i == 0 {
			// 如果job i是第1个job，那么只在前i个job中挑选的最优值为vi（vi为job i的价值）
			state[i][c][m] = State{
				Value: jobs[i].Value,
				Ids:   []int{jobs[i].Id},
			}
		} else {
			// 如果job i不是第1个job，那么只在前i个job中挑选的最优值根据状态转移方程更新
			nc := c - jobs[i].Cpu
			nm := m - jobs[i].Mem

			if state[i-1][c][m].Value > state[i-1][nc][nm].Value+jobs[i].Value {
				state[i][c][m] = state[i-1][c][m].DeepCopy()
			} else {
				state[i][c][m] = state[i-1][nc][nm].DeepCopy()
				state[i][c][m].Value += jobs[i].Value
				state[i][c][m].AddId(jobs[i].Id)
			}
		}
	}
}

func overResourceLimit(job Job, cpu int, mem int) bool {
	return job.Cpu > cpu || job.Mem > mem
}

func result() {
	c := int(cpuLimit)
	m := int(memLimit)
	fmt.Printf("%v --> %v\n", state[len(jobs)-1][c][m].Value, state[len(jobs)-1][c][m].Ids)
}

