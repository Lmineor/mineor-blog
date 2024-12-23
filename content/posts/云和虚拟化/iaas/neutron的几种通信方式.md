---
title: "一文说明白neutron用的几种通信方式"
date: 2024-12-19
draft: false
tags : [                    # 文章所属标签
    "Iaas",
]
---

# Neutron通信方式
neutron有如下通信方式：

- callback（同步调用）
- rpc(可以异步也可以同步)
- rest(同步调用)

## 1. callback

代码位置： `neutron-lib/neutron_lib/callbacks`

### 1.1 通信方式：

进程内的同步调用

### 1.2 原理

```py
def _get_callback_manager():
    global
    if _CALLBACK_MANAGER is None:
        _CALLBACK_MANAGER = manager.CallbacksManager()
    return _CALLBACK_MANAGER
```

通过`_get_callback_manager`初始化一个`_CALLBACK_MANAGER`（这里采用了python的单例模式），

再通过`registry.py`中的`subscribe`、`unsubscribe`、`unsubscribe_by_resource`、`unsubscribe_all`、`notify`、`publish`、`clear`函数来调用具体的回调函数。

下面看各个回调函数的具体的实现

#### 1.2.1 `subscribe`

订阅某个事件

```py
def subscribe(self, callback, resource, event,
              priority=priority_group.PRIORITY_DEFAULT):
    """Subscribe callback for a resource event.

    The same callback may register for more than one event.

    :param callback: 订阅的回调. It must raise or return a boolean.
    :param resource: 订阅的资源. It must be a valid resource.
    :param event: 订阅的事件. It must be a valid event.
    :param priority: the priority. Callbacks are sorted by priority
                      to be called. Smaller one is called earlier.
    """
    LOG.debug("Subscribe: %(callback)s %(resource)s %(event)s "
              "%(priority)d",
              {'callback': callback, 'resource': resource, 'event': event,
                'priority': priority})

    callback_id = _get_id(callback)

    # self._callbacks是个collections.defaultdict(dict)
    callbacks_list = self._callbacks[resource].setdefault(event, [])
    # 使用 setdefault 来确保 'event' 键存在，并为其设置默认值 []
    # 其中callbacks_list是个列表，若event存在，则返回已一直的列表，否则是个空列表
    # 格式如[(优先级,对应的callback),,,]
    # 对应callback的格式map key: callback_id, value callback
    for pc_pair in callbacks_list:
        # pc_pair格式 （优先级，对应的callback）
        if pc_pair[0] == priority:
            # 找到优先级和传参优先级相同callback，赋为pri_callbacks
            pri_callbacks = pc_pair[1]
            break
    else:
        # 没有找到优先级相同的callback
        pri_callbacks = {}
        # 按照格式，入callbacks_list
        callbacks_list.append((priority, pri_callbacks))
        callbacks_list.sort(key=lambda x: x[0])
    pri_callbacks[callback_id] = callback
    # 给callbacks_list添加当前订阅的callback

    # We keep a copy of callbacks to speed the unsubscribe operation.
    # 根据callback id简历所有，加速查找
    if callback_id not in self._index:
        self._index[callback_id] = collections.defaultdict(set)
    self._index[callback_id][resource].add(event)
```

总结： 这段代码主要用到了python的collections的defaultdict模块，简化了字典的操作。

#### 1.2.2 `unsubscribe`

取消某个事件的订阅

```py
def unsubscribe(self, callback, resource, event):
    """Unsubscribe callback from the registry.
    可以从事件级别删除订阅
    :param callback: 需要取消的订阅.
    :param resource: 需要取消订阅的资源.
    :param event: the event.
    """
    LOG.debug("Unsubscribe: %(callback)s %(resource)s %(event)s",
                {'callback': callback, 'resource': resource, 'event': event})

    callback_id = self._find(callback)
    if not callback_id:
        LOG.debug("Callback %s not found", callback_id)
        return
    if resource and event:
        # 先删除callback
        self._del_callback(self._callbacks[resource][event], callback_id)
        self._index[callback_id][resource].discard(event) # 先删除订阅的事件
        if not self._index[callback_id][resource]:
            # 该id的对应的订阅的资源的事件都被删掉了情况，则删除该资源
            del self._index[callback_id][resource]
            if not self._index[callback_id]:
                # 该id对应的订阅的资源都被删掉的情况，则删除该id
                del self._index[callback_id]
    else:
        value = '%s,%s' % (resource, event)
        raise exceptions.Invalid(element='resource,event', value=value)

def _find(self, callback):
    """Return the callback_id if found, None otherwise."""
    # 这里用到了上面提到的索引，可以根据id查看是否有该订阅
    callback_id = _get_id(callback)
    return callback_id if callback_id in self._index else None

def _del_callback(self, callbacks_list, callback_id):
    # callbacks_list格式如[(优先级,对应的callback),,,]
    # 对应callback的格式map key: callback_id, value callback
    for pc_pair in callbacks_list:
        # 通过for循环找到优先级和callback的元组
        pri_callbacks = pc_pair[1]
        # pri_callbacks是个字典：格式map key: callback_id, value callback
        if callback_id in pri_callbacks:
            del pri_callbacks[callback_id]
            if not pri_callbacks:
                # 若删除的是pri_callbacks里最后一个callback
                # 把callback list中的优先级和callback的元组删除掉，因为此时格式是这样的(优先级，{})，没有存在的必要了
                callbacks_list.remove(pc_pair)
            break
```

#### 1.2.3 `unsubscribe_by_resource`

取消某个资源的订阅

```py
def unsubscribe_by_resource(self, callback, resource):
    """Unsubscribe callback for any event associated to the resource.
    与unsubscribe大同小异
    :param callback: the callback.
    :param resource: the resource.
    """
    callback_id = self._find(callback)
    if callback_id:
        if resource in self._index[callback_id]:
            for event in self._index[callback_id][resource]:
                self._del_callback(self._callbacks[resource][event],
                                    callback_id)
            del self._index[callback_id][resource]
            if not self._index[callback_id]:
                del self._index[callback_id]
```


#### 1.2.4 `unsubscribe_all`

取消所有订阅

```py
def unsubscribe_all(self, callback):
    """Unsubscribe callback for all events and all resources.


    :param callback: the callback.
    """
    callback_id = self._find(callback)
    if callback_id:
        for resource, resource_events in self._index[callback_id].items():
            for event in resource_events:
                self._del_callback(self._callbacks[resource][event],
                                    callback_id)
        del self._index[callback_id]
```


#### 1.2.5 `notify`

给所有的订阅发消息(根据resource和event来区分)

```py
def notify(self, resource, event, trigger, **kwargs):
    """Notify all subscribed callback(s).

    Dispatch the resource's event to the subscribed callbacks.

    :param resource: The resource for the event.
    :param event: The event.
    :param trigger: The trigger. A reference to the sender of the event.
    :param kwargs: (deprecated) Unstructured key/value pairs to invoke
        the callback with. Using event objects with publish() is preferred.
    :raises CallbackFailure: CallbackFailure is raised if the underlying
        callback has errors.
    """
    # 一个个的通知
    errors = self._notify_loop(resource, event, trigger, **kwargs)
    if errors:
        # 处理错误的逻辑
        if event.startswith(events.BEFORE):
            abort_event = event.replace(
                events.BEFORE, events.ABORT)
            # 有before事件发生了错误，按照abort事件通知所有订阅的人，看是否需要做一个自己的操作（可以是回滚了什么的。。。）
            self._notify_loop(resource, abort_event, trigger, **kwargs)
            # 报个错
            raise exceptions.CallbackFailure(errors=errors)

        if event.startswith(events.PRECOMMIT):
            # precommit的直接报错(一般这种还没有入库，所以可以不用通知做其他事情)
            raise exceptions.CallbackFailure(errors=errors)

def _notify_loop(self, resource, event, trigger, **kwargs):
    """The notification loop."""
    # 开始一个个的通知
    errors = []
    # NOTE(yamahata): Since callback may unsubscribe it,
    # convert iterator to list to avoid runtime error.
    callbacks = list(itertools.chain(
        *[pri_callbacks.items() for (priority, pri_callbacks)
            in self._callbacks[resource].get(event, [])]))
    LOG.debug("Notify callbacks %s for %s, %s",
                [c[0] for c in callbacks], resource, event)
    # TODO(armax): consider using a GreenPile
    for callback_id, callback in callbacks:
        try:
            # 执行具体的调用了（你的callback就是一个函数）
            callback(resource, event, trigger, **kwargs)
        except Exception as e:
            abortable_event = (
                event.startswith(events.BEFORE) or
                event.startswith(events.PRECOMMIT)
            )
            if not abortable_event:
                LOG.exception("Error during notification for "
                                "%(callback)s %(resource)s, %(event)s",
                                {'callback': callback_id,
                                'resource': resource, 'event': event})
            else:
                LOG.debug("Callback %(callback)s raised %(error)s",
                            {'callback': callback_id, 'error': e})
            errors.append(exceptions.NotificationError(callback_id, e))
    return errors

```

#### 1.2.6 `publish`

发布订阅，作用和notify相同

```py
def publish(self, resource, event, trigger, payload=None):
    """Notify all subscribed callback(s) with a payload.

    Dispatch the resource's event to the subscribed callbacks.

    :param resource: The resource for the event.
    :param event: The event.
    :param trigger: The trigger. A reference to the sender of the event.
    :param payload: The optional event object to send to subscribers. If
        passed this must be an instance of BaseEvent.
    :raises neutron_lib.callbacks.exceptions.Invalid: if
        the payload object is not an instance of BaseEvent.
    :raises CallbackFailure: if the underlying callback has errors.
    """
    if payload:
        if not isinstance(payload, events.EventPayload):
            raise exceptions.Invalid(element='event payload',
                                        value=type(payload))
    return self.notify(resource, event, trigger, payload=payload)
```


#### 1.2.7 `clear`

把订阅相关的内存清楚掉

```py
def clear(self):
    """Brings the manager to a clean slate."""
    self._callbacks = collections.defaultdict(dict)
    self._index = collections.defaultdict(dict)
```

### 1.3 使用举例

```py
from neutron_lib.callbacks import registry

...

registry.subscribe(
    after_router_added, resources.ROUTER, events.AFTER_CREATE)

# after_router_added对应的callback
# resources.ROUTER 对应的资源
# events.AFTER_CREATE 对应的事件

def after_router_added(resource, event, l3_agent, **kwargs):
    router = kwargs['router']
    proxy = l3_agent.metadata_driver
    apply_metadata_nat_rules(router, proxy)
    if not isinstance(router, ha_router.HaRouter):
        proxy.spawn_monitored_metadata_proxy(
            l3_agent.process_monitor,
            router.ns_name,
            proxy.metadata_port,
            l3_agent.conf,
            router_id=router.router_id)
```

这样就完成了对router资源的，after_create事件的订阅，即当router的after_create事件发生，则会回调到after_router_added函数

其他事件也同理，就不举例了。

## 2. RPC

代码位置：oslo_messaging/rpc/client.py

neutron的RPC主要用到了oslo_messaging的功能，那么就主要讲讲oslo_messaging的rpc


据代码介绍，oslo_messaging的rpc支持两种模式，RPC calls 和 RPC casts.

- RPC calls 需要等待被调用方返回值

- RPC casts 不用等待被调用方返回值


先占个坑，以后再挖一下这个rpc 🤣🤣🤣

## 3. REST

这就不多说了，基操
