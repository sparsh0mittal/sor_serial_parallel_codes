//The use of this file is subject to the license terms in the LICENSE file 
//found in the top-level directory of this distribution. No part of this file,
// may be copied, modified, propagated, or distributed except according to 
//the terms contained in the LICENSE file.
//Author Sparsh Mittal sparsh0mittal@gmail.com



import std.stdio;
//import std.parallelism;
import std.concurrency;
import std.datetime;
import std.conv;
import std.math;
import core.sync.barrier;
__gshared Barrier myLocalBarrier = null;

immutable int gridSize = 4096;
immutable int gridSizePlus2 = gridSize+2;
immutable int MAXSTEPS = 5000_000;       /* Maximum number of iterations  */
immutable double TOL_VAL =0.00001;         /* Numerical Tolerance */
immutable double omega =  0.376;
immutable double one_minus_omega = 1.0 - omega;

immutable int NUM_CHECK =4000;

immutable int numberOfSlaves = 8;
immutable int factorValue = gridSize/numberOfSlaves;

double MAX_FUNC(double a, double b)
{
  return a> b? a: b;
}

double ABS_VAL(double a)
{
  return a> 0? a: -a;
}

shared double[gridSizePlus2][] gridInfo;
double[gridSizePlus2][] gridInfoOld;
double maxError = 0.0;

void main(string args[])
{
  writefln(" GridSize %s numberOfSlaves %s ", gridSize, numberOfSlaves);
   gridInfo = new shared double[gridSizePlus2][](gridSizePlus2);
   gridInfoOld = new double[gridSizePlus2][](gridSizePlus2);

  for(int i=0; i<gridSize+2; i++)
  {
    for(int j=0; j<gridSize+2; j++)
    {
      if(i==0)
        gridInfo[i][j] = 1.0;
      else
        gridInfo[i][j] = 0.0;
    }
  }

  bool shouldCheck = false;
  bool isConverged = false;
   StopWatch sw;
  sw.start(); //start/resume mesuring.
  for(int iter = 1; iter <= MAXSTEPS; iter++)
  {
    shouldCheck = false;
    if(iter % NUM_CHECK ==0)
    {
      shouldCheck = true;
      maxError = 0.0;

		for(int i=0; i<gridSize+2; i++)
		{
			for(int j=0; j<gridSize+2; j++)
			{			
				gridInfoOld[i][j] = gridInfo[i][j];			
			}
		}

    }


    {
      myLocalBarrier = new Barrier(numberOfSlaves+1);
      for (int cc=0; cc<numberOfSlaves; cc++) 
      {
         
	spawn(&SolverSlaveRed, thisTid,cc);          
      }
      
      //sync.
      myLocalBarrier.wait();
    }
    
    {
      myLocalBarrier = new Barrier(numberOfSlaves+1);
      for (int cc=0; cc<numberOfSlaves; cc++) 
      {
        
	 spawn(&SolverSlaveBlack, thisTid,cc);          
      } 
      //sync
      myLocalBarrier.wait();
    }

	if(shouldCheck)
     {	

		for(int i=0; i<gridSize+2; i++)
		{
			for(int j=0; j<gridSize+2; j++)
			{			
				maxError = fmax( abs(gridInfoOld[i][j] - gridInfo[i][j]), maxError);			
			}
		}

    if( maxError <  TOL_VAL)
      {
        isConverged = true;
        break;
      }
       else
	{
        writefln("Iter %s Error %s", iter, maxError);
        }
	}

  }
   sw.stop(); //stop/pause measuring.
  
  writeln(" Total time: ", (sw.peek().msecs/1000), "[sec]");
   if(isConverged)
    writeln("It converged");
  else
    writeln("It did not converge");
}



void SolverSlaveRed(Tid owner, int myNumber)
{
  
  double sum =0;
  int iStart = (myNumber*factorValue) + 1;
  int iEnd =  ((myNumber+1)*factorValue) ;
  
       
    //i even, j even
    for(int i=iStart+1; i<= iEnd; i+=2 )
    {
      for(int j=2; j< gridSize+1; j+=2)
      {
        
        
        
          sum = ( gridInfo[i  ][j+1] + gridInfo[i+1][j  ] +
              gridInfo[i-1][j  ] + gridInfo[i  ][j-1] )*0.25;
            
          
          gridInfo[i][j] = one_minus_omega*gridInfo[i][j] + omega*sum;          
              
          
      }
    }
    
    // i odd, j odd
    for(int i=iStart; i<= iEnd; i+=2)
    {
      for(int j=1; j< gridSize+1; j+=2)
      {
        
        
        
          sum = ( gridInfo[i  ][j+1] + gridInfo[i+1][j  ] +
              gridInfo[i-1][j  ] + gridInfo[i  ][j-1] )*0.25;
            
        
          gridInfo[i][j] = one_minus_omega*gridInfo[i][j] + omega*sum;          
              
          
      }
    }
  
  
  myLocalBarrier.wait();
 
  
  
  

    
  
}

void SolverSlaveBlack(Tid owner, int myNumber)
{
double sum =0;
  int iStart = (myNumber*factorValue) + 1;
  int iEnd =  ((myNumber+1)*factorValue) ;
  //i odd j even
    for(int i=iStart; i<= iEnd; i+=2)
      {
        for(int j=2; j< gridSize+1; j+=2)
        {
          
          
            sum = ( gridInfo[i  ][j+1] + gridInfo[i+1][j  ] +
                gridInfo[i-1][j  ] + gridInfo[i  ][j-1] )*0.25;
              
            
            gridInfo[i][j] = one_minus_omega*gridInfo[i][j] + omega*sum;          
                
            
        }
      }
      
      // i even j odd
      for(int i=iStart+1; i<= iEnd; i+=2)
      {
        for(int j=1; j< gridSize+1; j+=2)
        {
         
          
            sum = ( gridInfo[i  ][j+1] + gridInfo[i+1][j  ] +
                gridInfo[i-1][j  ] + gridInfo[i  ][j-1] )*0.25;
              
            
            gridInfo[i][j] = one_minus_omega*gridInfo[i][j] + omega*sum;          
                
           
        }
      }
  myLocalBarrier.wait();
}




