# Басплатное облако

## TL;DR;

Поднимаем инфраструктуру для хранения и сборки кода, а так же для деплоя личных сервисов. Для удобства и простоты код расположим в монорепозитории. Управление сервисов передадим в Managed Kubernetes, но с хранением конфигов в том же репозитории. Получившейся системой должно быть легко пользоваться, т.к. над ней есть полный контроль для любой настройки под себя.

## Требования

- Воспроизводимость. Все действия по созданию и настройке инфраструктуры должны быть описаны в документации или в скриптах. Это пригодится для восстановления контекста и для пересоздания облака заново.

- Минимимальная денежные затраты. Пока это не приносит денег, то и траты должны быть минимальными.

- Максимальная автоматизация. Действия, выполняемые во второй раз, нужно автоматизировать или обернуть в скрипт, которым легко пользоваться.

- Облачная сборка. Чтобы не мучить слабый компьютер и для возможности собирать по кнопке в браузерной IDE.

- Монорепозиторий.

## Содержимое облака

Т.к. в [Yandex Cloud](https://cloud.yandex.com) есть грант, то инфраструктуру будем строить в этом облаке. 

Для управления ресурсами используется [Managed Kubernetes](https://cloud.yandex.com/en-ru/services/managed-kubernetes), который для начала состоит из одной ноды и одного мастера. В случае нехватки ресурсов можно без проблем создать новую ноду с большим количеством CPU и Memory, в рамках доступного гранта.

Система сборки [Bazel](https://bazel.build) позволяет объединить весь код и действия над ним под одной системой сборки, позволяя использовать единый подход и инструментарий для сборки чего угодно. Также `Bazel` даёт возможность настроить удаленную сборку с кешированием, чтобы экономить ресурсы локальной машины и CI сервисов.

Для автоматической сборки, тестировани и выкладки используется [Drone CI](https://www.drone.io), который умеет интегрироваться с `GitHub` и выполняться любые действия по событиям, если обернуть их в `Docker container`.

## CI образ

CI образ нужен для того, чтобы собрать в одном месте, задокументировать и зафиксровать все зависимости, необходимые для сборки кода. Когда в репозиторий заливается новая версия кода, то для неё запускается pipeline сборки и тестирования, что по сути простые `bazel build \\...`
`bazel test \\...` команды. Чтобы не скачивать зависимости каждый раз, их можно сохранить в CI образ, который можно переиспользовать пока не они не меняются. Это позволяет экономить ресурсы и трафик при сборке и тестировании кода.

CI образ можно так же использовать для сборки и тестировани на локальной машине разработчика, на которой нет `bazel` или на нее не хочется его устанавливать. Для этого проднимается локальный контейнер, в который монтируется код. Локальная разработка через `CI образ` позволяет до PR убедиться, что сборка сборка герметична и никакая зависимость не пропущена.

## Кеширование

`Bazel` позволяет разбивает сбоку и тестирование на набор действий, упорядоченных в граф и кеширует результаты на каждом шаге, что значительно ускоряет повторный запуск команд. Например, если тест и тестируемый код не менялся после последнего успешного запуска, тогда еще раз тестировать не нужно. Более того, результаты промежуточных действий можно хранить в удаленном кеше и использовать на разных машинах. Т.е. после зеленого CI pipeline на локальной машине сборка завершится мгновенно, т.к. в `bazel` найдет в удаленном кеше успешные результаты сборки графа. 

## GitHub

Для удобства работы с кодом будем использвать GitHub, который надежнее чего-то, разавернутое внутри облака. Также такой выбор повышает надежность, т.к. при проблемых с собственной инфраструктурой код в GitHub будет доступен. В дальнейшен необходимо будет настроить зеркалирование репозитория в независимое место. 

## Выкладка через latest

В рамках CI pipeline не только тестируется код, но релизятся сервисы, что является тригером для перевыкладки. Для упрощения автоматической выкладки в kubernetes конфигах будет исользоваться latest, чтобы не приходилось его менять и руками выполнять деплой.