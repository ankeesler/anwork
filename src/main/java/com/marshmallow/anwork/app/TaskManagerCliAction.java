package com.marshmallow.anwork.app;

import com.marshmallow.anwork.app.cli.CliAction;
import com.marshmallow.anwork.core.FilePersister;
import com.marshmallow.anwork.core.Persister;
import com.marshmallow.anwork.task.TaskManager;

import java.io.File;
import java.util.Collection;
import java.util.Collections;

/**
 * This is a {@link CliAction} that loads a {@link TaskManager}, does something, and
 * then saves the {@link TaskManager}.
 *
 * @author Andrew
 * Created Sep 11, 2017
 */
public abstract class TaskManagerCliAction implements CliAction {

  private AnworkAppConfig config;

  public TaskManagerCliAction(AnworkAppConfig config) {
    this.config = config;
  }

  @Override
  public void run(String[] args) {
    try {
      TaskManager manager = loadTaskManager();
      run(args, manager);
      saveTaskManager(manager);
    } catch (Exception e) {
      System.out.println("Failed task manager action: " + e.getMessage());
    }
  }

  /**
   * Run the CLI action on a {@link TaskManager}.
   *
   * @param args The CLI arguments
   * @param manager The task manager
   */
  public abstract void run(String[] args, TaskManager manager);

  private TaskManager loadTaskManager() throws Exception {
    String context = config.getContext();
    File persistenceRoot = config.getPersistenceRoot();
    Persister<TaskManager> persister = new FilePersister<TaskManager>(persistenceRoot);
    if (!persister.contextExists(context)) {
      config.getDebugPrinter().accept("context " + context
                                      + " does not exist at root " + persistenceRoot + "!"
                                      + " creating new task manager");
      return new TaskManager();
    }

    Collection<TaskManager> loadeds = persister.load(context, TaskManager.serializer());
    if (loadeds.size() != 1) {
      throw new IllegalStateException("Persistence root " + persistenceRoot
                                      + " and context " + context
                                      + " contains 0 or more than 1 task manager.");
    }
    return loadeds.toArray(new TaskManager[0])[0];
  }

  private void saveTaskManager(TaskManager taskManager) throws Exception {
    Persister<TaskManager> persister = new FilePersister<TaskManager>(config.getPersistenceRoot());
    persister.save(config.getContext(),
                   TaskManager.serializer(),
                   Collections.singleton(taskManager));
  }
}
