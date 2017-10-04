package com.marshmallow.anwork.app;

import com.marshmallow.anwork.task.TaskManager;
import com.marshmallow.anwork.task.TaskState;

/**
 * This is a {@link TaskManagerCliAction} that sets the state of a task.
 *
 * <p>
 * Created Sep 11, 2017
 * </p>
 *
 * @author Andrew
 */
public class TaskManagerSetStateCliAction extends TaskManagerCliAction {

  private TaskState taskState;

  public TaskManagerSetStateCliAction(TaskState taskState) {
    super();
    this.taskState = taskState;
  }

  @Override
  public void run(AnworkAppConfig config, String[] args, TaskManager manager) {
    String taskName = args[0];
    manager.setState(taskName, taskState);
  }
}
