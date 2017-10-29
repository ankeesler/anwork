package com.marshmallow.anwork.app;

import com.marshmallow.anwork.app.cli.Action;
import com.marshmallow.anwork.app.cli.ActionCreator;
import com.marshmallow.anwork.app.cli.ArgumentType;
import com.marshmallow.anwork.app.cli.ArgumentValues;
import com.marshmallow.anwork.task.Task;
import com.marshmallow.anwork.task.TaskManager;
import com.marshmallow.anwork.task.TaskState;

/**
 * This is a {@link ActionCreator} for the task commands in the ANWORK app.
 *
 * <p>
 * Created Oct 4, 2017
 * </p>
 *
 * @author Andrew
 */
public class TaskCliActionCreator implements ActionCreator {

  private static class SetStateCliAction extends TaskManagerCliAction implements Action {

    private TaskState taskState;

    public SetStateCliAction(TaskState taskState) {
      super();
      this.taskState = taskState;
    }

    @Override
    public void run(AnworkAppConfig config,
                    ArgumentValues flags,
                    String[] args,
                    TaskManager manager) {
      String taskName = args[0];
      manager.setState(taskName, taskState);
    }
  }

  @Override
  public Action createAction(String commandName) {
    switch (commandName) {
      case "create":
        return new TaskManagerCliAction() {
          @Override
          public void run(AnworkAppConfig config,
                          ArgumentValues flags,
                          String[] args,
                          TaskManager manager) {
            String name = args[0];
            String description = (flags.containsKey("e")
                                  ? flags.getValue("e", ArgumentType.STRING)
                                  : "");
            int priority = (int)(flags.containsKey("p")
                                 ? flags.getValue("p", ArgumentType.NUMBER)
                                 : Task.DEFAULT_PRIORITY);
            manager.createTask(name, description, priority);
            config.getDebugPrinter().accept("created task '" + args[0] + "'");
          }
        };
      case "set-waiting":
        return new SetStateCliAction(TaskState.WAITING);
      case "set-blocked":
        return new SetStateCliAction(TaskState.BLOCKED);
      case "set-running":
        return new SetStateCliAction(TaskState.RUNNING);
      case "set-finished":
        return new SetStateCliAction(TaskState.FINISHED);
      case "delete":
        return new TaskManagerCliAction() {
          @Override
          public void run(AnworkAppConfig config,
                          ArgumentValues flags,
                          String[] args,
                          TaskManager manager) {
            manager.deleteTask(args[0]);
            config.getDebugPrinter().accept("deleted task '" + args[0] + "'");
          }
        };
      case "delete-all":
        return new TaskManagerCliAction() {
          @Override
          public void run(AnworkAppConfig config,
                          ArgumentValues flags,
                          String[] args,
                          TaskManager manager) {
            for (Task task : manager.getTasks()) {
              manager.deleteTask(task.getName());
            }
            config.getDebugPrinter().accept("deleted all tasks");
          }
        };
      case "show":
        return new TaskManagerCliAction() {
          @Override
          public void run(AnworkAppConfig config,
                          ArgumentValues flags,
                          String[] args,
                          TaskManager manager) {
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
