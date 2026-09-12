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
this value is written into separate result file

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
> more entry_list.txt 
10.
1.00563-9
17.28
17.28
17.28
4.9087390-6
.30
1.+9




READ_GROUNDING_FORCES
READ_MASSLESS_MECH
MPC_TO_CBUSH
GET_ID_RANGES

