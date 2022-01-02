package main

import (
"fmt"
)

func main() {
	// var a int
	// fmt.Scan(&a)
	fmt.Printf("%s", "hello world")
	nums:=[]int{1,3,5,2,4,6,8,0}
	res:=longestIncreasingSubArray(nums)
	fmt.Println(res)

}
func longestIncreasingSubArray(nums []int) []int{
	if len(nums)==0 {
		return []int{}
	}
	if len(nums)==1 {
		return nums
	}
	res :=[]int{}
	for i:=0;i<len(nums)-1;i++ {
		temp:=[]int{nums[i]}
		for j:=i+1;j<len(nums);j++ {
			if nums[j]>nums[j-1]{
				temp = append(temp, nums[j])
			}else {
				continue
			}
		}
		if len(res)<len(temp) {
			res=append([]int{},temp...)
		}
	}
	return res
}

