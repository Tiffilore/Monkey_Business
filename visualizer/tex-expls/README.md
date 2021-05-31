# Examples


File        | level     | verbosity | inclToken | goObjType | inclEnv   | input   
---         | ---       | ---       | ---       | ---       | ---       | --- 
ptree_0     | stmt      | 1         | false     | -         | -         | `let max = fn(x,y){if (x > y){return x} y}`
e_tree_0    | prog      | 0         | false     | false     | false     | `let a = 5 / 2 == 2`
e_tree_1    | prog      | 0         | false     | false     | true      | `let a = 5 / 2 == 2`