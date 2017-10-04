package com.marshmallow.anwork.app;

import com.marshmallow.anwork.app.cli.Cli;
import com.marshmallow.anwork.app.cli.CliAction;
import com.marshmallow.anwork.app.cli.CliList;
import com.marshmallow.anwork.journal.Journal;
import com.marshmallow.anwork.task.Task;
import com.marshmallow.anwork.task.TaskManager;
import com.marshmallow.anwork.task.TaskManagerJournalEntry;
import com.marshmallow.anwork.task.TaskState;

import java.util.Arrays;

/**
 * This class creates the CLI for the ANWORK app.
 *
 * <p>
 * Created Sep 11, 2017
 * </p>
 *
 * @author Andrew
 */
public class AnworkCliCreator {

  /**
   * Create an instance of the CLI for the ANWORK app.
   *
   * @return An instance of the CLI for the ANWORK app
   */
  public Cli makeCli() {
    Cli cli = new Cli("anwork", "ANWORK CLI commands");
    CliList root = cli.getRoot();
    makeRootFlags(root);
    makeTaskCommands(root);
    makeJournalCommands(root);
    return cli;
  }

  private void makeRootFlags(CliList root) {
    root.addLongFlag(AnworkAppCliFlag.DEBUG.getShortFlag(),
                     "debug",
                     "Turn on debug printing");
    root.addLongFlagWithParameter(AnworkAppCliFlag.CONTEXT.getShortFlag(),
                                  "context",
                                  "Set the persistence context",
                                  "name",
                                  "The name of the persistence context",
                                  AnworkAppCliFlag.CONTEXT.getParameterType());
    root.addLongFlagWithParameter(AnworkAppCliFlag.PERSISTENCE_ROOT.getShortFlag(),
                                  "output",
                                  "Set persistence output directory",
                                  "directory",
                                  "The directory at which to output the persistence data",
                                  AnworkAppCliFlag.PERSISTENCE_ROOT.getParameterType());
    root.addLongFlag(AnworkAppCliFlag.DONT_PERSIST.getShortFlag(),
                     "no-persist",
                     "Do not persist any task information");
  }

  private void makeTaskCommands(CliList root) {
    CliList taskCommandList = root.addList("task", "Task commands...");

    CliAction createAction = new TaskManagerCliAction() {
      @Override
      public void run(AnworkAppConfig config, String[] args, TaskManager manager) {
        manager.createTask(args[0], args[1], Integer.parseInt(args[2]));
        config.getDebugPrinter().accept("created task '" + args[0] + "'");
      }
    };
    taskCommandList.addCommand("create", "Create a task", createAction);

    taskCommandList.addCommand("set-waiting",
                               "Set a task as waiting",
                               new TaskManagerSetStateCliAction(TaskState.WAITING));
    taskCommandList.addCommand("set-blocked",
                               "Set a task as blocked",
                               new TaskManagerSetStateCliAction(TaskState.BLOCKED));
    taskCommandList.addCommand("set-running",
                               "Set a task as running",
                               new TaskManagerSetStateCliAction(TaskState.RUNNING));
    taskCommandList.addCommand("set-finished",
                               "Set a task as finished",
                               new TaskManagerSetStateCliAction(TaskState.FINISHED));

    CliAction deleteAction = new TaskManagerCliAction() {
      @Override
      public void run(AnworkAppConfig config, String[] args, TaskManager manager) {
        manager.deleteTask(args[0]);
        config.getDebugPrinter().accept("deleted task '" + args[0] + "'");
      }
    };
    taskCommandList.addCommand("delete", "Delete a task", deleteAction);

    CliAction showAction = new TaskManagerCliAction() {
      @Override
      public void run(AnworkAppConfig config, String[] args, TaskManager manager) {
        for (Task task : manager.getTasks()) {
          System.out.println(task);
        }
      }
    };
    taskCommandList.addCommand("show", "Show all tasks", showAction);
  }

  private void makeJournalCommands(CliList root) {
    CliList journalCommandList = root.addList("journal", "Journal commands...");

    CliAction showAllAction = new TaskManagerCliAction() {
      @Override
      public void run(AnworkAppConfig config, String[] args, TaskManager manager) {
        System.out.println(Arrays.toString(manager.getJournal().getEntries()));
      }
    };
    journalCommandList.addCommand("show-all",
                                  "Show all of the entries in the journal",
                                  showAllAction);

    CliAction showAction = new TaskManagerCliAction() {
      @Override
      public void run(AnworkAppConfig config, String[] args, TaskManager manager) {
        Journal<TaskManagerJournalEntry> journal = manager.getJournal(args[0]);
        if (journal == null) {
          System.out.println("No entries for task named " + args[0]);
        } else {
          System.out.println(Arrays.toString(journal.getEntries()));
        }
      }
    };
    journalCommandList.addCommand("show",
                                  "Show the entries in the journal for a task",
                                  showAction);
  }
}
