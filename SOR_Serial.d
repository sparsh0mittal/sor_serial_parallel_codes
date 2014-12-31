#!/usr/bin/env rdmd

//The use of this file is subject to the license terms in the LICENSE file 
//found in the top-level directory of this distribution. No part of this file,
// may be copied, modified, propagated, or distributed except according to 
//the terms contained in the LICENSE file.
//Author Sparsh Mittal sparsh0mittal@gmail.com



import std.stdio;

import std.algorithm;
import std.file;
import std.exception;
import std.string;
import std.concurrency;
import core.thread;
import std.datetime;
import std.process;
import std.conv;
import std.array;

import std.parallelism;

immutable int gridSize = 4096;
immutable double  INITIAL_GUESS = 0.0;
immutable int MAXSTEPS = 50000;       /* Maximum number of iterations               */
immutable double TOL_VAL =0.00001;         /* Numerical Tolerance */
immutable double  PI_VAL =3.14159265;       /* pi */
immutable int NUM_CHECK =4000;
immutable double omega =  0.376;
immutable double one_minus_omega = 1.0 - 0.376;

bool is_power_of_two(int x)
{
  return ( (x > 0) && ((x & (x - 1)) == 0) );
}

double MAX_FUNC(double a, double b)
{
  return a> b? a: b; 
}

double ABS_VAL(double a)
{
  return a> 0? a: -a;
}

void main(string args[])
{
  double maxError = 0.0;

  
  bool isConverged = false;
  
  string outputFileName =  "DResults/S"~"_G"~ to!string(gridSize+2);
  string timeFileName =  "DResults/TimeS"~"_G" ~ to!string(gridSize+2);
  
   
  double[gridSize+2][] gridInfo = new double[gridSize+2][](gridSize+2);
  
  for(int i=0; i<gridSize+2; i++)
  {
    for(int j=0; j<gridSize+2; j++)
    {
      if(i==0)
        gridInfo[i][j] = 1.0;
       else
        gridInfo[i][j] = INITIAL_GUESS;          
    }
  }
  StopWatch sw;
  sw.start(); //start/resume mesuring.
  //string phase = "red";
  double sum =0;
  for(int iter = 1; iter <= MAXSTEPS; iter++)
  {
    maxError = 0.0;
    
    /* process RED odd points ... */
    // i odd, j odd
    
    
    for(int i=1; i< gridSize+1; i+= 2)
    {
      for(int j=1; j< gridSize+1; j+= 2)
      {
        
            
        
        
          sum = ( gridInfo[i  ][j+1] + gridInfo[i+1][j  ] + gridInfo[i-1][j  ] + gridInfo[i  ][j-1] )*0.25;          
          
          maxError = MAX_FUNC(ABS_VAL(omega *(sum-gridInfo[i][j])), maxError);
          gridInfo[i][j] = one_minus_omega*gridInfo[i][j] + omega*sum;          
               
        
      }
    }
    
    
    /* process RED even points ... */
    // i even, j even
    for(int i=2; i< gridSize+1; i+= 2)
    {
      for(int j=2; j< gridSize+1; j+= 2)
      {
        
        
        
        sum = ( gridInfo[i  ][j+1] + gridInfo[i+1][j  ] + gridInfo[i-1][j  ] + gridInfo[i  ][j-1] )*0.25;          
          
        maxError = MAX_FUNC(ABS_VAL(omega *(sum-gridInfo[i][j])), maxError);
        gridInfo[i][j] = one_minus_omega*gridInfo[i][j] + omega*sum;          
               
        
      }
    }
    
    
    /* process BLACK odd points ... */
    // i odd, j even
    
    for(int i=1; i< gridSize+1; i+= 2)
    {
      for(int j=2; j< gridSize+1; j+= 2)
      {
        
        
          sum = ( gridInfo[i  ][j+1] + gridInfo[i+1][j  ] +
              gridInfo[i-1][j  ] + gridInfo[i  ][j-1] )*0.25;
          
          maxError = MAX_FUNC(ABS_VAL(omega* (sum-gridInfo[i][j])), maxError);
          gridInfo[i][j] = one_minus_omega*gridInfo[i][j] + omega*sum;        
               
        
      }
    }
    
    
    /* process BLACK even points ... */
    // i even, j odd
    
    for(int i=2; i< gridSize+1; i+= 2)
    {
      for(int j=1; j< gridSize+1; j+= 2)
      {
        
        
          sum = ( gridInfo[i  ][j+1] + gridInfo[i+1][j  ] +
              gridInfo[i-1][j  ] + gridInfo[i  ][j-1] )*0.25;
          
          maxError = MAX_FUNC(ABS_VAL(omega* (sum-gridInfo[i][j])), maxError);
          gridInfo[i][j] = one_minus_omega*gridInfo[i][j] + omega*sum;        
               
        
      }
    }
    
    if(iter % NUM_CHECK ==0)
    {
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
  
     auto timeHandle = File(timeFileName, "w");
     timeHandle.writeln( (sw.peek().msecs/1000));


  
}



