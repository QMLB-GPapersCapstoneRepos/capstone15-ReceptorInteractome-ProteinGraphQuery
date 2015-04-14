*****  DUE TO THE SERVER SPACE LIMIT, NOT ALL RECEPTOR'S PREDICTION FILES ARE HERE   ******

Totally there exist 905 prediction files. 
In this directory, we provide 318 of them (due to OUR SERVER HARD DISK SPACE LIMIT). 
The other files @ http://www.cs.cmu.edu/afs/cs.cmu.edu/project/structure-9/PPI/HMRI/allresults-eachReceptorOwn/

-----------------------------------------------

This directory contains the supplementary files for paper: 

Yanjun Qi[1]1, Harpreet K. Dhiman2, Neil Bhola3, Ivan Budyak4, Siddhartha Kar5, David Man2, Arpana Dutta2, Kalyan Tirupula2, Brian I. Carr5, Jennifer Grandis3,  Ziv Bar-Joseph1§ and Judith Klein-Seetharaman1,2,4§
Systematic prediction of human membrane receptor interactions, PROTEOMICS (2009) (In Press)
Supplementary Information see: http://www.cs.cmu.edu/~qyj/HMRI/

-----------------------------------------------
-----------------------------------------------

SubDir:  ./allresults-eachReceptorOwn

For each human membrane receptor, we predict its potential interaction partners from all human proteins. 
All possible pairs are tested and scores are provided (with no score cutoff). 

*****  DUE TO THE SERVER SPACE LIMIT, NOT ALL RECEPTOR'S PREDICTION FILES ARE HERE   ******

Totally there exist 905 prediction files. 
In this directory, we provide 318 of them (due to OUR SERVER HARD DISK SPACE LIMIT). 
The other files @ http://www.cs.cmu.edu/afs/cs.cmu.edu/project/structure-9/PPI/HMRI/allresults-eachReceptorOwn/


For each receptor, there exist a separate prediction file which is compressed into *.tar.gz format. 
You can use the gene ID of receptors to locate its prediction file in this subdir. 

Please use "gunzip" then "tar -xvf " to decompress it.  You can OPEN the file using EXCEL. 

In each line, Items are Separated by TAB.
The pairs could be sorted based on the "column E" .

column A, gene ID of gene1 ==> geneID1
column B, gene Symbol of gene1 ==> geneSym1
column C, gene ID of gene2 ==> geneID2
column D, gene Symbol of gene2 ==> geneSym2
column E, predicted score for this pair ==> RFpnScore
column F, the label of HPRD data set : [1 means interact; 2 means not-interact]  ==> hprdLabel
column H, gene description of gene1 ==> GeneID1:Synonyms:description
column J, features we used in the computational method ==> 27features
column L, GDB ID for gene1 ==> GDBID
column M, genetic disorder information for gene1 ==> GeneticDisorder

-----------------------------------------------
