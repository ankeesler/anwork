package com.marshmallow.anwork.app;

import com.marshmallow.anwork.app.cli.CliAction;
import com.marshmallow.anwork.app.cli.CliActionCreator;
import com.marshmallow.anwork.task.Task;
import com.marshmallow.anwork.task.TaskManager;
import com.marshmallow.anwork.task.TaskState;

/**
 * This is a {@link CliActionCreator} for the task commands in the ANWORK app.
 *
 * <p>
 * Created Oct 4, 2017
 * </p>
 *
 * @author Andrew
 */
public class TaskCliActionCreator implements CliActionCreator {

  @Override
  public CliAction createAction(String commandName) {
    switch (commandName) {
      case "create":
        return new TaskManagerCliAction() {
          @Override
          public void run(AnworkAppConfig config, String[] args, TaskManager manager) {
            manager.createTask(args[0], args[1], Integer.parseInt(args[2]));
            config.getDebugPrinter().accept("created task '" + args[0] + "'");
          }
        };
      case "set-waiting":
        return new TaskManagerSetStateCliAction(TaskState.WAITING);
      case "set-blocked":
        return new TaskManagerSetStateCliAction(TaskState.BLOCKED);
      case "set-running":
        return new TaskManagerSetStateCliAction(TaskState.RUNNING);
      case "set-finished":
        return new TaskManagerSetStateCliAction(TaskState.FINISHED);
      case "delete":
        return new TaskManagerCliAction() {
          @Override
          public void run(AnworkAppConfig config, String[] args, TaskManager manager) {
            manager.deleteTask(args[0]);
            config.getDebugPrinter().accept("deleted task '" + args[0] + "'");
          }
        };
      case "show":
        return new TaskManagerCliAction() {
          @Override
          public void run(AnworkAppConfig config, String[] args, TaskManager manager) {
            for (Task task : manager.getTasks()) {
              System.out.println(task);
            }
          }
        };
      default:
        return null; // error!
    }
  }
}
