package com.marshmallow.anwork.task;

/**
 * This is the state of a task.
 *
 * <p>
 * Created Aug 29, 2017
 * </p>
 *
 * @author Andrew
 */
public enum TaskState {
  WAITING,
  BLOCKED,
  RUNNING,
  FINISHED,
  ;
}
