# nas_tools
several tools processing MSC.Nastran files


actions:
--------
READ
STATS
SPLIT
EXTRACT_ACC_LIST
GET_CARD_ENTRY
READ_GROUNDING_FORCES
READ_MASSLESS_MECH
MPC_TO_CBUSH
GET_ID_RANGES







details for each action:
------------------------



READ

using arguments:
    "input_file": "grid.txt",
    "input_dir": "/home/lutz/prog/nas_tools/regression_tests"
    "output_file": "result.txt"

this utility reads the file gives line numbers and write the file into new file,
defined in output_file
can be used for debugging purposes and checking of base functionality

example output:
 INFO:  BEGIN BULK -  false
 lines/cards:  11268 11267 11267
 write nas cards into file: result.txt
    write:  result.txt




STATS

using arguments:
    "input_file": "grid.txt",
    "input_dir": "/home/lutz/prog/nas_tools/regression_tests"

gives back statistics in file as number of lines and cards and detailed card stats

can be used for debugging or first approach to a new file

example output:
INFO:  BEGIN BULK -  false
lines/cards:  84 51 51
... get stats 51
debug printout of nas card stats: 
DAREA     : 4
EIGRL     : 1
FREQ1     : 1
GRID      : 24
MDLPRM    : 2
PARAM     : 2
RBE2      : 16
RLOAD2    : 1
-----------
TOTAL     : 51
===========




SPLIT

using arguments:
    "input_file": "split_test_small.dat",
    "input_dir": "/home/lutz/prog/nas_tools/regression_tests"
    "output_dir": "result_dir"

splits the original input file into several separate files per card,
for example all GRID cards are splitted into grid.txt, all CBAR cards into cbar.txt

all splitted files are written into result_dir for easy later access to files

example output into result_dir:

INFO:  BEGIN BULK -  false
lines/cards:  84 51 51
created: result_dir/mdlprm.txt (2 cards)
created: result_dir/param.txt (2 cards)
created: result_dir/freq1.txt (1 cards)
created: result_dir/eigrl.txt (1 cards)
created: result_dir/grid.txt (24 cards)
created: result_dir/rbe2.txt (16 cards)
created: result_dir/darea.txt (4 cards)
created: result_dir/rload2.txt (1 cards)

example content of result_dir:
-rw-r--r-- 1 lutz lutz  164 Sep 12 17:58 darea.txt
-rw-r--r-- 1 lutz lutz   33 Sep 12 17:58 eigrl.txt
-rw-r--r-- 1 lutz lutz   40 Sep 12 17:58 freq1.txt
-rw-r--r-- 1 lutz lutz 1176 Sep 12 17:58 grid.txt
-rw-r--r-- 1 lutz lutz   36 Sep 12 17:58 mdlprm.txt
-rw-r--r-- 1 lutz lutz   29 Sep 12 17:58 param.txt
-rw-r--r-- 1 lutz lutz 2237 Sep 12 17:58 rbe2.txt
-rw-r--r-- 1 lutz lutz   49 Sep 12 17:58 rload2.txt







EXTRACT_ACC_LIST

using arguments:
 "input_file": "reg_test_01.dat",
 "input_dir": "/home/lutz/prog/nas_tools/regression_tests"
    
 "output_file": "result.txt",
 "output_dir": "result_dir",
 "option_01": "GRID",
 "input_01": "extract_ids.txt"

a list of certain entities according card_name and id can be extracted

an input file is scanned for cards with ids and the foundings are written 
into a separate file in result_dir/output_file. the option_01 defines the card_name,
the input_01 argument defines the required ids.

for example GRID entries with IDs 5,7,10 should be extracted, the GRID argument is
defined in option_01, the ids are stored in file defined by input_01 as list:
> more extract_ids.txt
5
7
10

terminal output:
INFO:  BEGIN BULK -  false
lines/cards:  65 21 21
extract cards ... GRID extract_ids.txt result.txt
num of IDs:  3

finally the output file contains cards according card_name and id:
> more result.txt 
GRID     5              .2       0.      0.
GRID     7              .3       0.      0.
GRID     10             .45      0.      0.




GET_CARD_ENTRY

get certain entry of a card and write this entry into list to file

so for example requesting each 5th entry of line 1 in PBEAM card,
this value is written into separate result file, example for PBEAM:

$  1   ||  2   ||  3   ||  4   ||  5   ||  6   ||  7   ||  8   ||  9   |
PBEAM   2       2       36.     17.28   675.            500.            


using arguments:
 "input_file": "reg_test_01.dat",
 "input_dir": "/home/lutz/prog/nas_tools/regression_tests"

 "output_file": "result.txt",
 "array_01": ["PBEAM",1,5]

terminal output:
INFO:  BEGIN BULK -  false
lines/cards:  65 21 21
PBEAM 1 5  -  5
entry list length:  8

output file containing requested values
> more result.txt 
10.
1.00563-9
17.28
17.28
17.28
4.9087390-6
.30
1.+9






READ_GROUNDING_FORCES

read grounding forces from .f06 according defined filter value

in .f06 file the grounding forces are a result of a failed groundcheck
and for each node with a grounding force this result could appear in .f06 file.
the grounding force itself contains 6 values for each DOF and the filter value can
be used to filter only grounding forces up to a certain level. 
for example if the filter value is set to 1.0E4, the maximum value of all six force
values is used as decider to filter that node as grounded node.

excerpt from .f06 file:
                                                                                     DIRECTION        6
                           G R O U N D   C H E C K   F O R C E S  ( F - S E T )

 POINT ID.   TYPE          T1             T2             T3             R1             R2             R3
 10670060      G      0.0            0.0            0.0            0.0            7.083236E+10   0.0
 10670070      G      0.0            0.0            0.0            0.0            8.452217E+10   0.0
 10670080      G      0.0            0.0            0.0            0.0            1.102255E+11   0.0
 10670090      G      7.229361E+10   0.0            0.0            0.0            1.074049E+11   6.051408E+10

using arguments:
  "input_file": "gf_and_mm.f06",
  "input_dir": "/home/lutz/prog/nas_tools/regression_tests"

  "output_file": "result.txt",
  "limit_size": 1.0E4
    
terminal output:
debug printout of control object: 
   Action:       READ_GROUNDING_FORCES
   InputFile:    gf_and_mm.f06
   InputDir:     /home/lutz/prog/nas_tools/regression_tests
   OutputFile:   result.txt
   OutputDir:    
   Option01:     
   Array01:      []
   Input01:      
   LimitSize:    10000
      FullInputPath:     /home/lutz/prog/nas_tools/regression_tests/gf_and_mm.f06
      FullOutpuPath:     result.txt
reading grounding forces ...
number of grounding nodes:  378

output files are result.ses and result.txt whereby result.ses contains the list 
of filtered nodes prepared for a patran session file that creates a group and 
result.txt as a list of nodes as simple text file:

==> result.ses <==
ga_group_create("grounding_nodes")
ga_group_entity_add("grounding_nodes", "node" // @ 
" 10670010 " // @ 
" 10670030 " // @ 
" 10670040 " // @ 
" 10670050 " // @ 
" 10670060 " // @ 
" 10670070 " // @ 
" 10670080 " // @ 
" 10670090 " // @ 

==> result.txt <==
10670010
10670030
10670040
10670050
10670060
10670070
10670080
10670090
10670530
10670540








READ_MASSLESS_MECH

get massless mechanism nodes from VAXW table in .f06 file

in Nastran runs a VAXW table can be printed containing nodes and DOFs with massless mechanisms, 
these nodes are filtered and will be output

excerpt from .f06:

0                VAXW
     POINT           VALUE     POINT           VALUE     POINT           VALUE     POINT           VALUE     POINT           VALUE


 COLUMN           1
   15201301 T1  1.00000E+00  15201301 T2  1.00000E+00  15201301 T3  1.00000E+00  15201301 R1  1.00000E+00  15201301 R2  1.00000E+00
   15201301 R3  1.00000E+00  15201302 T1  1.00000E+00  15201302 T2  1.00000E+00  15201302 T3  1.00000E+00  15201302 R1  1.00000E+00
   15201302 R2  1.00000E+00  15201302 R3  1.00000E+00  15201303 T1  1.00000E+00  15201303 T2  1.00000E+00  15201303 T3  1.00000E+00

using arguments:
   "input_file": "gf_and_mm.f06",
   "input_dir": "/home/lutz/prog/nas_tools/regression_tests"

   "output_file": "result.txt"

terminal output:
    
 debug printout of control object: 
   Action:       READ_MASSLESS_MECH
   InputFile:    gf_and_mm.f06
   InputDir:     /home/lutz/prog/nas_tools/regression_tests
   OutputFile:   result.txt
   OutputDir:    
   Option01:     
   Array01:      []
   Input01:      
   LimitSize:    0
      FullInputPath:     /home/lutz/prog/nas_tools/regression_tests/gf_and_mm.f06
      FullOutpuPath:     result.txt
reading massless mechanisms ...
number of massless nodes:  10   

output files are result.ses and result.txt whereby result.ses contains the list 
of filtered nodes prepared for a patran session file that creates a group and 
result.txt as a list of nodes as simple text file:

==> result.ses <==
ga_group_create("massless_nodes")
ga_group_entity_add("massless_nodes", "node" // @ 
" 15201301 " // @ 
" 15201302 " // @ 
" 15201303 " // @ 
" 15201304 " // @ 
" 15201306 " // @ 
" 15201307 " // @ 
" 15201309 " // @ 
" 15201310 " // @ 

==> result.txt <==
15201301
15201302
15201303
15201304
15201306
15201307
15201309
15201310
15201311
15201312






MPC_TO_CBUSH

Converts multiple MPC entries to CBUSH elements

To avoid possible grounding effects sometimes it is required to convert MPC entries with non-coincident
nodes to CBUSH elements which are real finite elements and which avoid possible grounding effects ( which
depends of course on the used DOF in MPC equations). MPC entries are collected and for each set of MPC
equations one CBUSH entry with same nodes is created. Element coordinate system is set to basic coordinate
system 0.

using arguments:
  "input_file": "mpc_small.txt",
  "input_dir": "/home/lutz/prog/nas_tools/regression_tests"

  "output_file": "result.txt"


terminal output:

debug printout of control object: 
   Action:       MPC_TO_CBUSH
   InputFile:    mpc_small.txt
   InputDir:     /home/lutz/prog/nas_tools/regression_tests
   OutputFile:   result.txt
   OutputDir:    
   Option01:     
   Array01:      []
   Input01:      
   LimitSize:    0
      FullInputPath:     /home/lutz/prog/nas_tools/regression_tests/mpc_small.txt
      FullOutpuPath:     result.txt
... read nastran cards:  /home/lutz/prog/nas_tools/regression_tests/mpc_small.txt
INFO:  BEGIN BULK -  false
lines/cards:  80 78 78
MPC_TO_CBUSH 50001 50001
lines:  78
   50000003 620620 6
   50000004 630620 3
   50000456 50000460 6
   50000457 50000461 6
   50000458 50000462 6
   50000459 50000463 6
   50000001 600620 6
   50000417 10800390 3
   50000379 50000417 6
   50000460 10800320 6
   50000461 10800330 6
   50000462 10800280 6
   50000463 10800290 6
   50000002 610620 6

for example input:

$  1   ||  2   ||  3   ||  4   ||  5   ||  6   ||  7   ||  8   ||  9   ||  10  |
MPC     9000000050000001       1     -1.  600620       1      1.
MPC     9000000050000001       2     -1.  600620       2      1.
MPC     9000000050000001       3     -1.  600620       3      1.
MPC     9000000050000001       4     -1.  600620       4      1.
MPC     9000000050000001       5     -1.  600620       5      1.
MPC     9000000050000001       6     -1.  600620       6      1.
MPC     9000000050000002       1     -1.  610620       1      1.
MPC     9000000050000002       2     -1.  610620       2      1.
MPC     9000000050000002       3     -1.  610620       3      1.

is converted to:

CBUSH   50001   50001   50000001600620                          0               
CBUSH   50002   50001   50000001600620                          0  








GET_ID_RANGES

get id range table for cards in input file

if looking for free id range an existing id range table can be created to get an overview
of already defined ids. the input file is read and an id range table is created for a certain range

using arguments:

  "input_file": "mpc_small.txt",
  "input_dir": "/home/lutz/prog/nas_tools/regression_tests"

  "output_dir": "result_dir"

terminal output:

debug printout of control object: 
   Action:       GET_ID_RANGES
   InputFile:    any_grids.dat
   InputDir:     /home/lutz/prog/nas_tools/regression_tests
   OutputFile:   
   OutputDir:    result_dir
   Option01:     
   Array01:      []
   Input01:      
   LimitSize:    0
      FullInputPath:     /home/lutz/prog/nas_tools/regression_tests/any_grids.dat
      FullOutpuPath:     result_dir
... read nastran cards:  /home/lutz/prog/nas_tools/regression_tests/any_grids.dat
INFO:  BEGIN BULK -  false
lines/cards:  24 24 24
     101      102      141      201      202      241      301      321      341      401 
     421      441    10000    10001    10002    10003    10004    37244  1010007  1010008 
 1010009  1010010  1010011  4000002 
length of id array:  24

example file contains ids:

GRID         101        -832.101-444.341353.2967
GRID         102        -1055.28-382.563  251.76
GRID         141        -825.291-563.402870.4579

id_file.txt is created:
0..10000 - 13  with last id: 10000
10000..20000 - 4  with last id: 10004
20000..30000 - 0
30000..40000 - 1  with last id: 37244
