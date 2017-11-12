package com.marshmallow.anwork.app;

import com.marshmallow.anwork.app.cli.Action;
import com.marshmallow.anwork.app.cli.Argument;
import com.marshmallow.anwork.app.cli.ArgumentType;
import com.marshmallow.anwork.app.cli.ArgumentValues;
import com.marshmallow.anwork.app.cli.Command;
import com.marshmallow.anwork.app.cli.Flag;
import com.marshmallow.anwork.core.FilePersister;
import com.marshmallow.anwork.core.Persister;
import com.marshmallow.anwork.task.Task;
import com.marshmallow.anwork.task.TaskManager;

import java.io.File;
import java.util.Collection;
import java.util.Collections;

/**
 * This is a {@link Action} that loads a {@link TaskManager}, does something, and
 * then saves the {@link TaskManager}.
 *
 * <p>
 * Created Sep 11, 2017
 * </p>
 *
 * @author Andrew
 */
public abstract class TaskManagerCliAction implements Action {

  /**
   * This is the name of the {@link Argument} used to specify the name of a {@link Task}. It is of
   * type {@link ArgumentType#STRING}.
   */
  protected static final String TASK_SPECIFIER_ARGUMENT = "task-specifier";

  @Override
  public void run(ArgumentValues flags, ArgumentValues arguments) {
    AnworkAppConfig config = new AnworkAppConfig(flags);
    try {
      TaskManager manager = loadTaskManager(config);
      if (run(config, flags, arguments, manager)) {
        saveTaskManager(config, manager);
      }
    } catch (Exception e) {
      throw new IllegalStateException("Failed task manager action!", e);
    }
  }

  /**
   * Run the CLI action on a {@link TaskManager}.
   *
   * @param config The {@link AnworkAppConfig} applied to this {@link Action}
   * @param flags The {@link Flag} {@link ArgumentValues} passed to this {@link Command}
   * @param arguments The {@link ArgumentValues} passed to this {@link Command}
   * @param manager The task manager
   * @return Whether or not the {@code manager} should be persisted to file.
   */
  public abstract boolean run(AnworkAppConfig config,
                              ArgumentValues flags,
                              ArgumentValues arguments,
                              TaskManager manager);

  protected static String getTaskSpecifierArgument(TaskManager manager, ArgumentValues arguments) {
    TaskSpecifierParser taskSpecifierParser = new TaskSpecifierParser(manager);
    String taskSpecifier = arguments.getValue(TASK_SPECIFIER_ARGUMENT, ArgumentType.STRING);

    // See the javadoc for TaskSpecifierArgumentType#convert for the reason why we know this
    // variable is non-null and not empty.
    return taskSpecifierParser.parse(taskSpecifier)[0];
  }

  protected static Persister<TaskManager> getPersister(AnworkAppConfig config) {
    File persistenceRoot = config.getPersistenceRoot();
    Persister<TaskManager> persister = new FilePersister<TaskManager>(persistenceRoot);
    return persister;
  }

  private TaskManager loadTaskManager(AnworkAppConfig config) throws Exception {
    String context = config.getContext();
    File persistenceRoot = config.getPersistenceRoot();
    Persister<TaskManager> persister = getPersister(config);
    if (!config.getDoPersist()) {
      config.getDebugPrinter().accept("not loading task manager because"
                                      + " persist command line option set to false");
      return new TaskManager();
    } else if (!persister.exists(context)) {
      config.getDebugPrinter().accept("context " + context
                                      + " does not exist at root " + persistenceRoot + "!"
                                      + " creating new task manager");
      return new TaskManager();
    }

    Collection<TaskManager> loadeds = persister.load(context, TaskManager.SERIALIZER);
    if (loadeds.size() != 1) {
      throw new IllegalStateException("Persistence root " + persistenceRoot
                                      + " and context " + context
                                      + " contains 0 or more than 1 task manager.");
    }
    return loadeds.toArray(new TaskManager[0])[0];
  }

  private void saveTaskManager(AnworkAppConfig config, TaskManager taskManager) throws Exception {
    Persister<TaskManager> persister = getPersister(config);
    if (!config.getDoPersist()) {
      config.getDebugPrinter().accept("not saving task manager because"
                                      + " no-persist command line option set");
    } else {
      persister.save(config.getContext(),
                     TaskManager.SERIALIZER,
                     Collections.singleton(taskManager));
    }
  }
}
